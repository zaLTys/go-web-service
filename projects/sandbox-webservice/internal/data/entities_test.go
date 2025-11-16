package data

import "testing"

func TestEntityRepositoryInterfaceCompliance(t *testing.T) {
	var _ EntityRepository = EntityModel{} // compile-time check
}
