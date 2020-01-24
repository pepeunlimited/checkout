// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/pepeunlimited/checkout/internal/app/app1/ent/vault"
)

// VaultCreate is the builder for creating a Vault entity.
type VaultCreate struct {
	config
	amount           *int64
	reference_number *string
}

// SetAmount sets the amount field.
func (vc *VaultCreate) SetAmount(i int64) *VaultCreate {
	vc.amount = &i
	return vc
}

// SetReferenceNumber sets the reference_number field.
func (vc *VaultCreate) SetReferenceNumber(s string) *VaultCreate {
	vc.reference_number = &s
	return vc
}

// Save creates the Vault in the database.
func (vc *VaultCreate) Save(ctx context.Context) (*Vault, error) {
	if vc.amount == nil {
		return nil, errors.New("ent: missing required field \"amount\"")
	}
	if vc.reference_number == nil {
		return nil, errors.New("ent: missing required field \"reference_number\"")
	}
	if err := vault.ReferenceNumberValidator(*vc.reference_number); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"reference_number\": %v", err)
	}
	return vc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (vc *VaultCreate) SaveX(ctx context.Context) *Vault {
	v, err := vc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (vc *VaultCreate) sqlSave(ctx context.Context) (*Vault, error) {
	var (
		v     = &Vault{config: vc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: vault.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: vault.FieldID,
			},
		}
	)
	if value := vc.amount; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  *value,
			Column: vault.FieldAmount,
		})
		v.Amount = *value
	}
	if value := vc.reference_number; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: vault.FieldReferenceNumber,
		})
		v.ReferenceNumber = *value
	}
	if err := sqlgraph.CreateNode(ctx, vc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	v.ID = int(id)
	return v, nil
}