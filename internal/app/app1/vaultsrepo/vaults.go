package vaultsrepo

import (
	"context"
	"errors"
	"github.com/pepeunlimited/checkout/internal/app/app1/ent"
	"github.com/pepeunlimited/checkout/internal/app/app1/ent/vault"
	"log"
)

var (
	ErrReferenceNumberExist = errors.New("vaults: reference number exist")
	ErrVaultsNotExist = errors.New("vaults: not exist")
)

type VaultsRepository interface {
	Add(ctx context.Context, amount int64, referenceNumber string) (*ent.Vault, *ent.Tx, error)
	Sum(ctx context.Context) 														(int64, error)
	Delete(ctx context.Context)
	GetByReferenceNumber(ctx context.Context, referenceNumber string)				(*ent.Vault, error)
}

type vaultsMySQL struct {
	client *ent.Client
}

func (v vaultsMySQL) GetByReferenceNumber(ctx context.Context, referenceNumber string) (*ent.Vault, error) {
	vault, err := v.client.Vault.Query().Where(vault.ReferenceNumber(referenceNumber)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrVaultsNotExist
		}
		return nil, err
	}
	return vault, nil
}

func (v vaultsMySQL) Sum(ctx context.Context) (int64, error) {
	panic("implement me")
}

func (v vaultsMySQL) Delete(ctx context.Context) {
	v.client.Vault.Delete().ExecX(ctx)
}

func (v vaultsMySQL) Add(ctx context.Context, amount int64, referenceNumber string) (*ent.Vault, *ent.Tx, error) {
	tx, err := v.client.Tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	vaults, err := tx.Vault.Create().SetAmount(amount).SetReferenceNumber(referenceNumber).Save(ctx)
	if err != nil {
		rollback(tx)
		if ent.IsConstraintError(err) {
			return nil, nil, ErrReferenceNumberExist
		}
		return nil, nil, err
	}
	return vaults, tx, nil
}

func rollback(tx *ent.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Print("vaults-repository: rollback failed: "+err.Error())
	}
}

func NewVaultsRepository(client *ent.Client) VaultsRepository {
	return vaultsMySQL{client:client}
}