package section

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNewImage(t *testing.T) {
	code := []uint64{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
	}
	cs := NewCodeSection(code)
	i := NewImage()
	err := i.AddSection(cs, true)
	if err != nil {
		t.Error(err)
		return
	}
	err = i.AddSection(cs, true)
	if err == nil {
		t.Errorf("expected already seen")
		return
	}
}

func TestVariable(t *testing.T) {
	v := Variable{
		Id:    1000,
		Name:  "Moo",
		value: new(big.Rat).SetFloat64(3.0),
	}
	ve, err := encodeVariableElement(&v)
	if err != nil {
		t.Error(err)
		return
	}
	vd, err := decodeVariableElement(ve, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(v, vd) {
		// expected due to big.Rat
		r1 := v.value.(*big.Rat)
		r2 := vd.value.(*big.Rat)
		if !reflect.DeepEqual(*r1, *r2) {
			t.Errorf("variables not equal")
			t.Logf("%v%v", *r1, *r2)
		}
		return
	}
}

func TestNewUncompressed(t *testing.T) {
	image := []uint64{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
	}
	cs := NewCodeSection(image)
	raw, err := cs.Raw(false)
	if err != nil {
		t.Error(err)
		return
	}

	dcs, err := SectionFromImage(raw, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(cs, dcs) {
		t.Errorf("uncompressed test corrupt")
		return
	}
}

func TestNewCompressedCodeSection(t *testing.T) {
	image := []uint64{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
	}
	cs := NewCodeSection(image)
	raw, err := cs.Raw(true)
	if err != nil {
		t.Error(err)
		return
	}

	dcs, err := SectionFromImage(raw, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(cs, dcs) {
		t.Errorf("uncompressed test corrupt")
		return
	}
}

func TestNewCompressedVariableSection(t *testing.T) {
	v1, err := NewVariable(1000, "Moo", new(big.Rat).SetFloat64(3.0))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = encodeVariableElement(v1)
	if err != nil {
		t.Error(err)
		return
	}
	v2, err := NewVariable(1001, "NotMoo", new(big.Rat).SetFloat64(5.0))
	_, err = encodeVariableElement(v2)
	if err != nil {
		t.Error(err)
		return
	}

	vs, err := NewVariableSection([]*Variable{v1, v2})
	if err != nil {
		t.Error(err)
		return
	}

	raw, err := vs.Raw(true)
	if err != nil {
		t.Error(err)
		return
	}

	dvs, err := SectionFromImage(raw, nil)
	if err != nil {
		t.Error(err)
		return
	}
	dvs = dvs // shut up
	// this does not test the variables
	//if !reflect.DeepEqual(vs, dvs) {
	//	t.Errorf("uncompressed test corrupt")
	//	t.Logf("%v%v", spew.Sdump(vs), spew.Sdump(dvs))
	//	return
	//}
}

func TestMulti(t *testing.T) {
	// vars
	v1, err := NewVariable(1000, "Moo", new(big.Rat).SetFloat64(3.0))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = encodeVariableElement(v1)
	if err != nil {
		t.Error(err)
		return
	}
	v2, err := NewVariable(1001, "NotMoo", new(big.Rat).SetFloat64(5.0))
	_, err = encodeVariableElement(v2)
	if err != nil {
		t.Error(err)
		return
	}

	vs, err := NewVariableSection([]*Variable{v1, v2})
	if err != nil {
		t.Error(err)
		return
	}

	// code
	code := []uint64{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
	}
	cs := NewCodeSection(code)

	// image
	i := NewImage()
	err = i.AddSection(cs, true)
	if err != nil {
		t.Error(err)
		return
	}
	err = i.AddSection(vs, true)
	if err != nil {
		t.Error(err)
		return
	}

	// decode from image
	sections, err := SectionsFromImage(i.GetImage())
	if err != nil {
		t.Error(err)
		return
	}

	// add corruption test here
	sections = sections // shut compiler up for now
}
