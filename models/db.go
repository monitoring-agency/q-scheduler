package models

import "github.com/myOmikron/echotools/utilitymodels"

type About struct {
	utilitymodels.CommonID
	Version string
}
