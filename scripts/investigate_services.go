package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/siderolabs/omni/client/pkg/client"
)

// This script investigates the Omni client services to document available methods
func main() {
	// Note: This requires a valid Omni client connection
	// Run with: go run scripts/investigate_services.go
	
	fmt.Println("=== Omni Client Service Investigation ===\n")
	
	// We can't actually create a client without credentials, but we can document
	// what we need to investigate
	
	fmt.Println("Services to investigate:")
	fmt.Println("1. Management Service (client.Management())")
	fmt.Println("2. Talos Service (client.Talos())")
	fmt.Println("3. Auth Service (client.Auth())")
	fmt.Println("4. OIDC Service (client.OIDC())\n")
	
	fmt.Println("Investigation approach:")
	fmt.Println("- Use reflection to inspect method signatures")
	fmt.Println("- Use go doc to view package documentation")
	fmt.Println("- Examine source code if available")
	fmt.Println("- Look for example usage in tests or documentation\n")
	
	// Example of how to inspect a service
	fmt.Println("Example inspection pattern:")
	fmt.Println(`
	mgmtClient := client.Management()
	clientType := reflect.TypeOf(mgmtClient)
	
	// List all methods
	for i := 0; i < clientType.NumMethod(); i++ {
		method := clientType.Method(i)
		fmt.Printf("Method: %s\n", method.Name)
		fmt.Printf("  Signature: %s\n", method.Type)
	}
	`)
	
	// Document what we're looking for
	fmt.Println("\n=== What to Document ===")
	fmt.Println("\nFor each service, document:")
	fmt.Println("1. Client struct type")
	fmt.Println("2. All public methods")
	fmt.Println("3. Method signatures (parameters and return types)")
	fmt.Println("4. Request/Response types")
	fmt.Println("5. Error handling patterns")
	fmt.Println("6. Context usage")
	fmt.Println("7. Common patterns")
}

// Helper function to inspect a service client
func inspectService(service interface{}, serviceName string) {
	fmt.Printf("\n=== %s Service ===\n", serviceName)
	
	serviceType := reflect.TypeOf(service)
	if serviceType == nil {
		fmt.Printf("Service is nil or not initialized\n")
		return
	}
	
	fmt.Printf("Type: %s\n", serviceType.String())
	fmt.Printf("Kind: %s\n", serviceType.Kind())
	
	if serviceType.Kind() == reflect.Ptr {
		elemType := serviceType.Elem()
		fmt.Printf("Element Type: %s\n", elemType.String())
		
		// List methods
		fmt.Printf("\nMethods:\n")
		for i := 0; i < serviceType.NumMethod(); i++ {
			method := serviceType.Method(i)
			if !strings.HasPrefix(method.Name, "Unsafe") {
				fmt.Printf("  - %s%s\n", method.Name, formatMethodSignature(method.Type))
			}
		}
	}
}

func formatMethodSignature(methodType reflect.Type) string {
	if methodType.Kind() != reflect.Func {
		return ""
	}
	
	var params []string
	for i := 0; i < methodType.NumIn(); i++ {
		param := methodType.In(i)
		params = append(params, param.String())
	}
	
	var returns []string
	for i := 0; i < methodType.NumOut(); i++ {
		ret := methodType.Out(i)
		returns = append(returns, ret.String())
	}
	
	sig := "(" + strings.Join(params, ", ") + ")"
	if len(returns) > 0 {
		sig += " (" + strings.Join(returns, ", ") + ")"
	}
	
	return sig
}
