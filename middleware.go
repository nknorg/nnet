package nnet

import (
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/util"
)

// ApplyMiddleware add a middleware to node, network, router, etc. If multiple
// middleware of the same type are applied, they will be called in the order of
// being added.
func (nn *NNet) ApplyMiddleware(mw interface{}) error {
	applied := false
	errs := util.NewErrors()

	err := nn.GetLocalNode().ApplyMiddleware(mw)
	if err == nil {
		applied = true
	} else {
		errs = append(errs, err)
	}

	err = nn.Network.ApplyMiddleware(mw)
	if err == nil {
		applied = true
	} else {
		errs = append(errs, err)
	}

	for _, router := range nn.GetRouters() {
		err = router.ApplyMiddleware(mw)
		if err == nil {
			applied = true
		} else {
			errs = append(errs, err)
		}
	}

	if !applied {
		return errs.Merged()
	}

	return nil
}

// MustApplyMiddleware is the same as ApplyMiddleware, but will panic if an
// error occurs. This is a convenient shortcut if ApplyMiddleware is not
// expected to fail.
func (nn *NNet) MustApplyMiddleware(mw interface{}) {
	err := nn.ApplyMiddleware(mw)
	if err != nil {
		log.Errorf("Apply middleware error: %v", err)
		panic(err)
	}
}
