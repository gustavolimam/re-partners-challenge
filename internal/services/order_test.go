package services

import (
	"re-partners-challenge/internal/clients"
	"re-partners-challenge/internal/constants"
	"re-partners-challenge/internal/models"
	"reflect"
	"testing"
)

func TestOrder_CalculateOrderPackQty(t *testing.T) {
	type fields struct {
		cache func() *clients.Cache
	}
	type args struct {
		order models.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.OrderPacks
		wantErr bool
	}{
		{
			name: "success_case_1",
			fields: fields{
				cache: func() *clients.Cache {
					cache := clients.NewCache()
					cache.Set(constants.PackSizesCacheKey, constants.PackSizesDefault, 5)
					return cache
				},
			},
			args: args{
				order: models.Order{
					Items: 1500,
				},
			},
			want: []models.OrderPacks{
				{
					Size:  1000,
					Count: 1,
				},
				{
					Size:  500,
					Count: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "success_case_2",
			fields: fields{
				cache: func() *clients.Cache {
					cache := clients.NewCache()
					cache.Set(constants.PackSizesCacheKey, constants.PackSizesDefault, 5)
					return cache
				},
			},
			args: args{
				order: models.Order{
					Items: 6750,
				},
			},
			want: []models.OrderPacks{
				{
					Size:  5000,
					Count: 1,
				},
				{
					Size:  1000,
					Count: 1,
				},
				{
					Size:  500,
					Count: 1,
				},
				{
					Size:  250,
					Count: 1,
				},
			},
		},
		{
			name: "success_case_3",
			fields: fields{
				cache: func() *clients.Cache {
					cache := clients.NewCache()
					cache.Set(constants.PackSizesCacheKey, []int{4, 100}, 0)
					return cache
				},
			},
			args: args{
				order: models.Order{
					Items: 10,
				},
			},
			want: []models.OrderPacks{
				{
					Size:  4,
					Count: 3,
				},
			},
		},
		{
			name: "error_no_items_to_pack",
			fields: fields{
				cache: func() *clients.Cache {
					return clients.NewCache()
				},
			},
			args: args{
				order: models.Order{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error_cache",
			fields: fields{
				cache: func() *clients.Cache {
					cache := clients.NewCache()
					return cache
				},
			},
			args: args{
				order: models.Order{
					Items: 1500,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.fields.cache()

			o := &Order{
				cache: cache,
			}
			got, err := o.CalculateOrderPackQty(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateOrderPackQty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateOrderPackQty() got = %v, want %v", got, tt.want)
			}
		})
	}
}
