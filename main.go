package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/omni/client/pkg/client"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"

	omniclient "github.com/jubblin/omni-api/internal/client"
	"github.com/jubblin/omni-api/internal/i18n"
	resconverter "github.com/jubblin/omni-api/internal/resource"
)

// Version is set at build time via ldflags
var Version = "dev"

const (
	clusterNodePrefixMachineSets     = "machinesets-"
	clusterNodePrefixClusterMachines = "clustermachines-"
	labelKeyMachineSet               = "omni.sidero.dev/machine-set"
)

type ResourceQuery struct {
	Type   resource.Type
	IsList bool
}

type TreeNode struct {
	ID       widget.TreeNodeID
	Type     string
	Label    string
	Resource map[string]interface{}
	Children []*TreeNode
}

type AppState struct {
	stateClient              state.State
	selectedResourceType     string
	treeRoot                 *TreeNode
	resourceTree             *widget.Tree
	resourceIDInput          *widget.Entry
	statusLabel              *widget.Label
	detailTitle              *widget.Label
	detailText               *widget.RichText
	currentResource          map[string]interface{}
	machineLinksContainer    *fyne.Container
	versionLinksContainer    *fyne.Container
	machineSetLinksContainer *fyne.Container
	machineRelatedLinksContainer *fyne.Container
}

type DetailComponents struct {
	Title              *widget.Label
	Text               *widget.RichText
	MachineLinks       *fyne.Container
	VersionLinks       *fyne.Container
	MachineSetLinks    *fyne.Container
	MachineRelatedLinks *fyne.Container
}

type ResourceContext struct {
	AppState    *AppState
	StateClient state.State
	Ctx         context.Context
	Depth       int
}

func getResourceQueries() map[string]ResourceQuery {
	return map[string]ResourceQuery{
		i18n.T("resource.cluster"): {
			Type:   omni.ClusterType,
			IsList: true,
		},
		i18n.T("resource.cluster_machines"): {
			Type:   omni.ClusterMachineType,
			IsList: true,
		},
		i18n.T("resource.etcd_backups"): {
			Type:   omni.EtcdBackupType,
			IsList: true,
		},
		i18n.T("resource.kubernetes_versions"): {
			Type:   omni.KubernetesVersionType,
			IsList: true,
		},
		i18n.T("resource.machine"): {
			Type:   omni.MachineType,
			IsList: true,
		},
		i18n.T("resource.machine_classes"): {
			Type:   omni.MachineClassType,
			IsList: true,
		},
		i18n.T("resource.machine_set"): {
			Type:   omni.MachineSetType,
			IsList: true,
		},
		i18n.T("resource.ongoing_tasks"): {
			Type:   omni.OngoingTaskType,
			IsList: true,
		},
		i18n.T("resource.schematics"): {
			Type:   omni.SchematicType,
			IsList: true,
		},
	}
}

func initializeApp(myWindow fyne.Window) (*client.Client, *AppState) {
	omniClient, err := omniclient.NewOmniClient()
	if err != nil {
		log.Printf("Failed to create Omni client: %v", err)
		// Format directive is in translation file: "app.error.client" = "Error: %v\n\nPlease set OMNI_ENDPOINT..."
		errorLabel := widget.NewLabel(i18n.T("app.error.client", err)) //nolint
		errorLabel.Wrapping = fyne.TextWrapWord
		myWindow.SetContent(container.NewScroll(errorLabel))
		return nil, nil
	}

	appState := &AppState{
		stateClient: omniClient.Omni().State(),
		treeRoot: &TreeNode{
			ID:       "",
			Type:     "root",
			Label:    i18n.T("tree.root"),
			Children: []*TreeNode{},
		},
	}

	return omniClient, appState
}

func createResourceIDInput() *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(i18n.T("resource.id.placeholder"))
	return entry
}

func createResourceTree(appState *AppState) *widget.Tree {
	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			return getNodeChildren(id, appState)
		},
		func(id widget.TreeNodeID) bool {
			return isBranch(id, appState)
		},
		func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			node := getNodeByID(id, appState)
			if node != nil {
				label.SetText(node.Label)
			}
		},
	)
	return tree
}

func getNodeChildren(id widget.TreeNodeID, appState *AppState) []widget.TreeNodeID {
	if id == "" {
		if appState.treeRoot == nil {
			return []widget.TreeNodeID{}
		}
		children := make([]widget.TreeNodeID, 0, len(appState.treeRoot.Children))
		for _, child := range appState.treeRoot.Children {
			children = append(children, child.ID)
		}
		return children
	}

	node := getNodeByID(id, appState)
	if node == nil {
		return []widget.TreeNodeID{}
	}

	children := make([]widget.TreeNodeID, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, child.ID)
	}
	return children
}

func isBranch(id widget.TreeNodeID, appState *AppState) bool {
	node := getNodeByID(id, appState)
	if node == nil {
		return false
	}
	return len(node.Children) > 0
}

func getNodeByID(id widget.TreeNodeID, appState *AppState) *TreeNode {
	if id == "" {
		return appState.treeRoot
	}
	return findNodeRecursive(appState.treeRoot, id)
}

func findNodeRecursive(node *TreeNode, id widget.TreeNodeID) *TreeNode {
	if node.ID == id {
		return node
	}
	for _, child := range node.Children {
		if found := findNodeRecursive(child, id); found != nil {
			return found
		}
	}
	return nil
}

func setupTreeSelection(resourceTree *widget.Tree, appState *AppState) {
	resourceTree.OnSelected = func(id widget.TreeNodeID) {
		node := getNodeByID(id, appState)
		if node != nil && node.Resource != nil {
			updateDetailPane(node.Resource, appState)
		}
	}
}

func createDetailComponents() *DetailComponents {
	title := widget.NewLabel(i18n.T("tree.select_resource"))
	title.TextStyle = fyne.TextStyle{Bold: true}

	text := widget.NewRichText()
	text.Wrapping = fyne.TextWrapWord

	return &DetailComponents{
		Title:               title,
		Text:                text,
		MachineLinks:        container.NewVBox(),
		VersionLinks:        container.NewVBox(),
		MachineSetLinks:     container.NewVBox(),
		MachineRelatedLinks: container.NewVBox(),
	}
}

func setupAppState(appState *AppState, resourceTree *widget.Tree, resourceIDInput *widget.Entry, statusLabel *widget.Label, detailComponents *DetailComponents) {
	appState.detailText = detailComponents.Text
	appState.detailTitle = detailComponents.Title
	appState.resourceTree = resourceTree
	appState.resourceIDInput = resourceIDInput
	appState.statusLabel = statusLabel
	appState.machineLinksContainer = detailComponents.MachineLinks
	appState.versionLinksContainer = detailComponents.VersionLinks
	appState.machineSetLinksContainer = detailComponents.MachineSetLinks
	appState.machineRelatedLinksContainer = detailComponents.MachineRelatedLinks
}

func createBurgerMenu(myWindow fyne.Window, appState *AppState, executeQuery func()) *widget.Button {
	resourceQueries := getResourceQueries()
	resourceOptions := make([]string, 0, len(resourceQueries))
	for name := range resourceQueries {
		resourceOptions = append(resourceOptions, name)
	}
	sort.Strings(resourceOptions)

	menuItems := make([]*fyne.MenuItem, 0, len(resourceOptions))
	for _, name := range resourceOptions {
		name := name
		menuItems = append(menuItems, fyne.NewMenuItem(name, func() {
			appState.selectedResourceType = name
			executeQuery()
		}))
	}

	menu := fyne.NewMenu(i18n.T("resource.type.menu"), menuItems...)
	burgerMenu := widget.NewButton(fmt.Sprintf("☰ %s", appState.selectedResourceType), nil)
	
	burgerMenu.OnTapped = func() {
		popup := widget.NewPopUpMenu(menu, fyne.CurrentApp().Driver().CanvasForObject(burgerMenu))
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(burgerMenu)
		popup.ShowAtPosition(pos.Add(fyne.NewPos(0, burgerMenu.Size().Height)))
	}

	return burgerMenu
}

func createLanguageSelector(appState *AppState, refreshUI func()) *widget.Select {
	languages := i18n.GetSupportedLanguages()
	options := make([]string, 0, len(languages))
	for _, lang := range languages {
		options = append(options, i18n.GetLanguageName(lang))
	}

	langSelect := widget.NewSelect(options, func(selected string) {
		for _, lang := range languages {
			if i18n.GetLanguageName(lang) == selected {
				if err := i18n.SetLanguage(lang); err == nil {
					refreshUI()
				}
				break
			}
		}
	})

	currentLang := i18n.GetCurrentLanguage()
	langSelect.SetSelected(i18n.GetLanguageName(currentLang))

	return langSelect
}

func createMainLayout(burgerMenu *widget.Button, resourceIDInput *widget.Entry, resourceTree *widget.Tree, detailComponents *DetailComponents, statusLabel *widget.Label, langSelect *widget.Select) *container.Split {
	topBar := container.NewBorder(nil, nil, nil, langSelect, burgerMenu)
	middleBar := container.NewBorder(nil, nil, widget.NewLabel(i18n.T("resource.id.label")), nil, resourceIDInput)
	// Format directive is in translation file: "app.connected" = "Connected to: %s"
	infoLabel := widget.NewLabel(i18n.T("app.connected", os.Getenv("OMNI_ENDPOINT"))) //nolint
	infoLabel.Wrapping = fyne.TextWrapWord

	leftPane := container.NewBorder(topBar, middleBar, nil, nil, resourceTree)
	
	detailScroll := container.NewScroll(container.NewVBox(
		detailComponents.Title,
		detailComponents.Text,
		detailComponents.MachineLinks,
		detailComponents.VersionLinks,
		detailComponents.MachineSetLinks,
		detailComponents.MachineRelatedLinks,
	))
	detailScroll.SetMinSize(fyne.NewSize(400, 0))

	rightPane := container.NewBorder(nil, statusLabel, nil, nil, detailScroll)
	
	split := container.NewHSplit(leftPane, rightPane)
	split.SetOffset(0.3)

	return split
}

func createExecuteQueryFunc(appState *AppState, resourceIDInput *widget.Entry, statusLabel *widget.Label, stateClient state.State) func() {
	return func() {
		resourceQueries := getResourceQueries()
		query, ok := resourceQueries[appState.selectedResourceType]
		if !ok {
			statusLabel.SetText(i18n.T("resource.type.unknown"))
			return
		}

		statusLabel.SetText(i18n.T("resource.querying"))
		ctx := context.Background()

		if query.IsList {
			executeListQuery(ctx, query.Type, appState, stateClient, statusLabel)
		} else {
			resourceID := resourceIDInput.Text
			if resourceID == "" {
				statusLabel.SetText(i18n.T("resource.id.required"))
				return
			}
			executeGetQuery(ctx, query.Type, resourceID, appState, stateClient, statusLabel)
		}
	}
}

func executeListQuery(ctx context.Context, resourceType resource.Type, appState *AppState, stateClient state.State, statusLabel *widget.Label) {
	md := resource.NewMetadata(omniresources.DefaultNamespace, resourceType, "", resource.VersionUndefined)
	list, err := stateClient.List(ctx, md)
	if err != nil {
		handleQueryError(err, appState, statusLabel)
		return
	}

	resources := make([]map[string]interface{}, 0, len(list.Items))
	for _, item := range list.Items {
		resources = append(resources, resconverter.ToMap(item))
	}

	appState.treeRoot.Children = buildResourceNodes(resources, stateClient, ctx)
	appState.resourceTree.Refresh()
	statusLabel.SetText(i18n.T("resource.success", fmt.Sprintf("%d resources", len(resources))))
}

func executeGetQuery(ctx context.Context, resourceType resource.Type, resourceID string, appState *AppState, stateClient state.State, statusLabel *widget.Label) {
	md := resource.NewMetadata(omniresources.DefaultNamespace, resourceType, resourceID, resource.VersionUndefined)
	item, err := stateClient.Get(ctx, md)
	if err != nil {
		handleQueryError(err, appState, statusLabel)
		return
	}

	resourceMap := resconverter.ToMap(item)
	appState.treeRoot.Children = []*TreeNode{
		{
			ID:       fmt.Sprintf("%s-%s", resourceType, resourceID),
			Type:     string(resourceType),
			Label:    formatResourceLabel(resourceMap),
			Resource: resourceMap,
			Children: []*TreeNode{},
		},
	}
	appState.resourceTree.Refresh()
	statusLabel.SetText(i18n.T("resource.success", resourceID))
}

func handleQueryError(err error, appState *AppState, statusLabel *widget.Label) {
	// Format directive is in translation file: "resource.error" = "Error: %v"
	statusLabel.SetText(i18n.T("resource.error", err)) //nolint
	appState.treeRoot.Children = []*TreeNode{}
	appState.resourceTree.Refresh()
}

func setupResourceIDAutoQuery(resourceIDInput *widget.Entry, appState *AppState, executeQuery func()) {
	resourceIDInput.OnSubmitted = func(_ string) {
		executeQuery()
	}
}

func buildResourceNodes(resources []map[string]interface{}, stateClient state.State, ctx context.Context) []*TreeNode {
	nodes := make([]*TreeNode, 0, len(resources))
	for _, resourceMap := range resources {
		resourceID, resourceType := extractResourceInfoFromMap(resourceMap)
		enrichMachineResource(resourceMap, resourceType, resourceID, stateClient, ctx, 0)
		enrichClusterMachineResource(resourceMap, resourceType, resourceID, stateClient, ctx, 0)
		node := createResourceNodeFromMap(resourceMap, resourceID, resourceType)
		nodes = append(nodes, node)
	}
	return nodes
}

func extractResourceInfoFromMap(resourceMap map[string]interface{}) (string, string) {
	resourceID := ""
	if id, ok := resourceMap["id"].(string); ok {
		resourceID = id
	}
	resourceType := ""
	if t, ok := resourceMap["type"].(string); ok {
		resourceType = t
	}
	return resourceID, resourceType
}

func enrichMachineResource(resourceMap map[string]interface{}, resourceType, resourceID string, stateClient state.State, ctx context.Context, depth int) {
	if resourceType != string(omni.MachineType) || resourceID == "" || depth > 1 {
		return
	}

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineStatusType, resourceID, resource.VersionUndefined)
	machineStatus, err := stateClient.Get(ctx, md)
	if err != nil {
		return
	}

	ms, ok := machineStatus.(*omni.MachineStatus)
	if !ok {
		return
	}

	if ms.TypedSpec().Value.Network != nil && ms.TypedSpec().Value.Network.Hostname != "" {
		resourceMap["hostname"] = ms.TypedSpec().Value.Network.Hostname
	}
}

func enrichClusterMachineResource(resourceMap map[string]interface{}, resourceType, resourceID string, stateClient state.State, ctx context.Context, depth int) {
	if resourceType != string(omni.ClusterMachineType) || resourceID == "" || depth > 1 {
		return
	}

	machineID, ok := resourceMap["machine_id"].(string)
	if !ok || machineID == "" {
		return
	}

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineStatusType, machineID, resource.VersionUndefined)
	machineStatus, err := stateClient.Get(ctx, md)
	if err != nil {
		return
	}

	ms, ok := machineStatus.(*omni.MachineStatus)
	if !ok {
		return
	}

	if ms.TypedSpec().Value.Network != nil && ms.TypedSpec().Value.Network.Hostname != "" {
		resourceMap["hostname"] = ms.TypedSpec().Value.Network.Hostname
	}
}

func createResourceNodeFromMap(resourceMap map[string]interface{}, resourceID, resourceType string) *TreeNode {
	label := formatResourceLabel(resourceMap)
	return &TreeNode{
		ID:       fmt.Sprintf("%s-%s", resourceType, resourceID),
		Type:     resourceType,
		Label:    label,
		Resource: resourceMap,
		Children: []*TreeNode{},
	}
}

func formatResourceLabel(resourceMap map[string]interface{}) string {
	resourceID := ""
	if id, ok := resourceMap["id"].(string); ok {
		resourceID = id
	}
	resourceType := ""
	if t, ok := resourceMap["type"].(string); ok {
		resourceType = t
	}

	label := resourceID
	switch resourceType {
	case string(omni.ClusterType):
		if kv, ok := resourceMap["kubernetes_version"].(string); ok {
			label = fmt.Sprintf("%s (K8s: %s)", resourceID, kv)
		}
	case string(omni.MachineType):
		if hostname, ok := resourceMap["hostname"].(string); ok && hostname != "" {
			label = hostname
		} else if addr, ok := resourceMap["management_address"].(string); ok && addr != "" {
			label = fmt.Sprintf("%s (%s)", resourceID, addr)
		}
	case string(omni.ClusterMachineType):
		if hostname, ok := resourceMap["hostname"].(string); ok && hostname != "" {
			label = hostname
		} else if machineID, ok := resourceMap["machine_id"].(string); ok && machineID != "" {
			label = fmt.Sprintf("%s (Machine: %s)", resourceID, machineID)
		}
	case string(omni.MachineSetType):
		if mc, ok := resourceMap["machine_class"].(string); ok {
			label = fmt.Sprintf("%s (Class: %s)", resourceID, mc)
		}
	}
	return label
}

func updateDetailPane(resourceData map[string]interface{}, appState *AppState) {
	appState.currentResource = resourceData

	resourceID, resourceType := extractResourceInfo(resourceData)
	updateDetailTitle(appState, resourceID)
	displayResource := loadMachineStatusIfNeeded(resourceData, resourceType, resourceID, appState)
	updateDetailJSON(displayResource, appState)
	updateDetailLinks(resourceData, resourceType, resourceID, appState)
}

func extractResourceInfo(resourceData map[string]interface{}) (string, string) {
	resourceID := ""
	if id, ok := resourceData["id"].(string); ok {
		resourceID = id
	}
	resourceType := ""
	if t, ok := resourceData["type"].(string); ok {
		resourceType = t
	}
	return resourceID, resourceType
}

func updateDetailTitle(appState *AppState, resourceID string) {
	if resourceID != "" {
		appState.detailTitle.SetText(fmt.Sprintf("Resource Details: %s", resourceID))
	} else {
		appState.detailTitle.SetText("Resource Details")
	}
}

func loadMachineStatusIfNeeded(resourceData map[string]interface{}, resourceType, resourceID string, appState *AppState) map[string]interface{} {
	if resourceType != string(omni.MachineType) || resourceID == "" {
		return resourceData
	}

	ctx := context.Background()
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineStatusType, resourceID, resource.VersionUndefined)
	machineStatus, err := appState.stateClient.Get(ctx, md)
	if err != nil {
		return resourceData
	}

	statusMap := resconverter.ToMap(machineStatus)
	if statusSpec, ok := statusMap["spec"].(map[string]interface{}); ok {
		resourceData["machine_status"] = statusSpec
	}

	return resourceData
}

func updateDetailJSON(resourceData map[string]interface{}, appState *AppState) {
	jsonBytes, err := json.MarshalIndent(resourceData, "", "  ")
	if err != nil {
		appState.detailText.ParseMarkdown(fmt.Sprintf("Error formatting JSON: %v", err))
		return
	}
	appState.detailText.ParseMarkdown(fmt.Sprintf("```json\n%s\n```", string(jsonBytes)))
}

func updateDetailLinks(resourceData map[string]interface{}, resourceType, resourceID string, appState *AppState) {
	clearAllLinkContainers(appState)
	updateMachineIDLinks(resourceData, appState)
	updateKubernetesVersionLinks(resourceData, appState)
	updateMachineSetLinks(resourceData, appState)
}

func clearAllLinkContainers(appState *AppState) {
	appState.machineLinksContainer.RemoveAll()
	appState.versionLinksContainer.RemoveAll()
	appState.machineSetLinksContainer.RemoveAll()
	appState.machineRelatedLinksContainer.RemoveAll()
}

func updateMachineIDLinks(resourceData map[string]interface{}, appState *AppState) {
	machineIDs := findMachineIDs(resourceData)
	if len(machineIDs) == 0 {
		return
	}

	linksLabel := widget.NewLabel(i18n.T("detail.machine_id.label"))
	appState.machineLinksContainer.Add(linksLabel)
	for _, machineID := range machineIDs {
		// Format directive is in translation file: "detail.machine_id" = "Machine ID: %s"
		btn := widget.NewButton(i18n.T("detail.machine_id", machineID), func(id string) func() { //nolint
			return func() {
				loadMachineByID(id, appState)
			}
		}(machineID))
		appState.machineLinksContainer.Add(btn)
	}
}

func findMachineIDs(resourceData map[string]interface{}) []string {
	var machineIDs []string
	if machineID, ok := resourceData["machine_id"].(string); ok && machineID != "" {
		machineIDs = append(machineIDs, machineID)
	}
	return machineIDs
}

func updateKubernetesVersionLinks(resourceData map[string]interface{}, appState *AppState) {
	k8sVersions := findKubernetesVersions(resourceData)
	if len(k8sVersions) == 0 {
		return
	}

	linksLabel := widget.NewLabel(i18n.T("detail.kubernetes_version.label"))
	appState.versionLinksContainer.Add(linksLabel)
	for _, version := range k8sVersions {
		// Format directive is in translation file: "detail.kubernetes_version" = "Kubernetes Version: %s"
		btn := widget.NewButton(i18n.T("detail.kubernetes_version", version), func(v string) func() { //nolint
			return func() {
				loadResourcesByK8sVersion(v, appState)
			}
		}(version))
		appState.versionLinksContainer.Add(btn)
	}
}

func findKubernetesVersions(resourceData map[string]interface{}) []string {
	var versions []string
	if version, ok := resourceData["kubernetes_version"].(string); ok && version != "" {
		versions = append(versions, version)
	}
	return versions
}

func updateMachineSetLinks(resourceData map[string]interface{}, appState *AppState) {
	machineSetIDs := findMachineSetIDs(resourceData)
	if len(machineSetIDs) == 0 {
		return
	}

	linksLabel := widget.NewLabel(i18n.T("detail.machine_set.label"))
	appState.machineSetLinksContainer.Add(linksLabel)
	for _, machineSetID := range machineSetIDs {
		// Format directive is in translation file: "detail.machine_set" = "Machine Set: %s"
		btn := widget.NewButton(i18n.T("detail.machine_set", machineSetID), func(id string) func() { //nolint
			return func() {
				loadMachinesByMachineSet(id, appState)
			}
		}(machineSetID))
		appState.machineSetLinksContainer.Add(btn)
	}
}

func findMachineSetIDs(resourceData map[string]interface{}) []string {
	var machineSetIDs []string
	if labels, ok := resourceData["labels"].(map[string]interface{}); ok {
		if machineSet, ok := labels[labelKeyMachineSet].(string); ok && machineSet != "" {
			machineSetIDs = append(machineSetIDs, machineSet)
		}
	}
	return machineSetIDs
}

func loadMachineByID(machineID string, appState *AppState) {
	ctx := context.Background()
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineType, machineID, resource.VersionUndefined)

	machine, err := appState.stateClient.Get(ctx, md)
	if err != nil {
		// Format directives are in translation file: "machine.error.extensions" = "Error loading extensions for machine %s: %v"
		appState.statusLabel.SetText(i18n.T("machine.error.extensions", machineID, err)) //nolint
		return
	}

	machineMap := resconverter.ToMap(machine)
	updateDetailPane(machineMap, appState)
	// Format directive is in translation file: "machine.loaded" = "Loaded machine: %s"
	appState.statusLabel.SetText(i18n.T("machine.loaded", machineID)) //nolint
}

func loadMachinesByMachineSet(machineSetID string, appState *AppState) {
	ctx := context.Background()
	resources := findClusterMachinesByMachineSet(ctx, machineSetID, appState)

	if len(resources) == 0 {
		// Format directive is in translation file: "machineset.not_found" = "No machines found in machine set %s"
		appState.statusLabel.SetText(i18n.T("machineset.not_found", machineSetID)) //nolint
		return
	}

	appState.treeRoot.Children = buildResourceNodes(resources, appState.stateClient, ctx)
	appState.resourceTree.Refresh()
	// Format directives are in translation file: "machineset.found" = "Found %d machines in machine set %s"
	appState.statusLabel.SetText(i18n.T("machineset.found", len(resources), machineSetID)) //nolint
}

func findClusterMachinesByMachineSet(ctx context.Context, machineSetID string, appState *AppState) []map[string]interface{} {
	resources := make([]map[string]interface{}, 0)
	clusterMachineMD := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterMachineType, "", resource.VersionUndefined)
	list, err := appState.stateClient.List(ctx, clusterMachineMD)
	if err != nil {
		return resources
	}

	for _, item := range list.Items {
		cm, ok := item.(*omni.ClusterMachine)
		if !ok {
			continue
		}
		labels := cm.Metadata().Labels()
		if labels == nil {
			continue
		}
		machineSet, ok := labels.Get(labelKeyMachineSet)
		if !ok || machineSet != machineSetID {
			continue
		}
		resources = append(resources, resconverter.ToMap(cm))
	}

	return resources
}

func loadResourcesByK8sVersion(version string, appState *AppState) {
	ctx := context.Background()
	resources := findResourcesByK8sVersion(ctx, version, appState)

	if len(resources) == 0 {
		// Format directive is in translation file: "k8s.version.not_found" = "No resources found with Kubernetes version %s"
		appState.statusLabel.SetText(i18n.T("k8s.version.not_found", version)) //nolint
		return
	}

	appState.treeRoot.Children = buildResourceNodes(resources, appState.stateClient, ctx)
	appState.resourceTree.Refresh()
	// Format directives are in translation file: "k8s.version.found" = "Found %d resources with Kubernetes version %s"
	appState.statusLabel.SetText(i18n.T("k8s.version.found", len(resources), version)) //nolint
}

func findResourcesByK8sVersion(ctx context.Context, version string, appState *AppState) []map[string]interface{} {
	resources := make([]map[string]interface{}, 0)
	resources = append(resources, findClustersByK8sVersion(ctx, version, appState)...)
	resources = append(resources, findClusterMachinesByK8sVersion(ctx, version, appState)...)
	return resources
}

func findClustersByK8sVersion(ctx context.Context, version string, appState *AppState) []map[string]interface{} {
	resources := make([]map[string]interface{}, 0)
	clusterMD := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterType, "", resource.VersionUndefined)
	list, err := appState.stateClient.List(ctx, clusterMD)
	if err != nil {
		return resources
	}

	for _, item := range list.Items {
		c, ok := item.(*omni.Cluster)
		if !ok {
			continue
		}
		if c.TypedSpec().Value.KubernetesVersion == version {
			resources = append(resources, resconverter.ToMap(c))
		}
	}

	return resources
}

func findClusterMachinesByK8sVersion(ctx context.Context, version string, appState *AppState) []map[string]interface{} {
	resources := make([]map[string]interface{}, 0)
	clusterMachineMD := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterMachineType, "", resource.VersionUndefined)
	list, err := appState.stateClient.List(ctx, clusterMachineMD)
	if err != nil {
		return resources
	}

	for _, item := range list.Items {
		cm, ok := item.(*omni.ClusterMachine)
		if !ok {
			continue
		}
		if cm.TypedSpec().Value.KubernetesVersion == version {
			resources = append(resources, resconverter.ToMap(cm))
		}
	}

	return resources
}

func main() {
	if err := i18n.Init("en"); err != nil {
		log.Printf("Failed to initialize i18n: %v", err)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow(i18n.T("app.title"))
	myWindow.Resize(fyne.NewSize(1400, 900))

	omniClient, appState := initializeApp(myWindow)
	if omniClient == nil {
		return
	}
	defer omniClient.Close()

	appState.selectedResourceType = i18n.T("resource.cluster")

	resourceIDInput := createResourceIDInput()
	resourceTree := createResourceTree(appState)
	setupTreeSelection(resourceTree, appState)
	detailComponents := createDetailComponents()
	statusLabel := widget.NewLabel(i18n.T("app.ready"))
	executeQuery := createExecuteQueryFunc(appState, resourceIDInput, statusLabel, omniClient.Omni().State())

	setupResourceIDAutoQuery(resourceIDInput, appState, executeQuery)
	burgerMenu := createBurgerMenu(myWindow, appState, executeQuery)

	setupAppState(appState, resourceTree, resourceIDInput, statusLabel, detailComponents)

	refreshUI := func() {
		resourceQueries := getResourceQueries()
		resourceOptions := make([]string, 0, len(resourceQueries))
		for name := range resourceQueries {
			resourceOptions = append(resourceOptions, name)
		}
		sort.Strings(resourceOptions)
		burgerMenu.SetText(fmt.Sprintf("☰ %s", appState.selectedResourceType))

		resourceIDInput.SetPlaceHolder(i18n.T("resource.id.placeholder"))
		statusLabel.SetText(i18n.T("app.ready"))
		appState.treeRoot.Label = i18n.T("tree.root")
		if appState.detailTitle != nil {
			appState.detailTitle.SetText(i18n.T("tree.select_resource"))
		}
		appState.resourceTree.Refresh()
	}

	langSelect := createLanguageSelector(appState, refreshUI)
	mainContent := createMainLayout(burgerMenu, resourceIDInput, resourceTree, detailComponents, statusLabel, langSelect)

	myWindow.SetContent(mainContent)

	executeQuery()
	myWindow.ShowAndRun()
}
