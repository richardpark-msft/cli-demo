package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Auth struct {
	namespace string
	cs        string
}

func AddAuth(fs *flag.FlagSet) *Auth {
	auth := &Auth{}

	fs.StringVar(&auth.namespace, "ns", "", "Namespace (assumes DefaultAzureCredential)")
	fs.StringVar(&auth.cs, "csvar", "SERVICEBUS_CONNECTION_STRING", "Environment variable that contains a connection string")

	return auth
}

func (a *Auth) NewClient() (*azservicebus.Client, error) {
	var client *azservicebus.Client

	if a.namespace != "" {
		dac, err := azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			return nil, err
		}

		client, err = azservicebus.NewClient(a.namespace, dac, nil)

		if err != nil {
			return nil, err
		}
	} else {
		cs := os.Getenv(a.cs)

		if cs == "" {
			return nil, fmt.Errorf("no connection string in environment variable %s", cs)
		}

		var err error
		client, err = azservicebus.NewClientFromConnectionString(cs, nil)

		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
