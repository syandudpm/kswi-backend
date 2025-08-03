package oss

import (
	"context"
	"kswi-backend/internal/shared/pagination"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	DtDatabase(ctx context.Context, req DtDatabaseRequest) ([]DtDatabaseResponse, int, int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) DtDatabase(ctx context.Context, req DtDatabaseRequest) ([]DtDatabaseResponse, int, int, error) {
	var data []DtDatabaseResponse
	var total int
	var totalFiltered int
	var err error
	var total64 int64

	query := r.db.WithContext(ctx).Table("kswi.oss_base")

	// First get total count without any filters
	if err := query.Count(&total64).Error; err != nil {
		return nil, 0, 0, err
	}
	total = int(total64)

	// Initialize where conditions and parameters
	whereAnd := []string{}
	paramsAnd := []interface{}{}
	whereOr := []string{}
	paramsOr := []interface{}{}

	// Always include this condition
	// whereAnd = append(whereAnd, "statusPM IS NOT NULL")

	// Handle date filters
	if req.StartDate != nil && !req.StartDate.IsZero() {
		whereAnd = append(whereAnd, "_created_at >= ?")
		paramsAnd = append(paramsAnd, req.StartDate)
	}

	if req.EndDate != nil && !req.EndDate.IsZero() {
		whereAnd = append(whereAnd, "_created_at <= ?")
		paramsAnd = append(paramsAnd, req.EndDate)
	}

	// Handle JSON filters if they exist
	if req.Filters != nil {
		// Process AND filters
		for _, filter := range req.Filters.And {
			whereClause, param := pagination.BuildWhereClause(filter)
			if whereClause != "" {
				whereAnd = append(whereAnd, whereClause)
				paramsAnd = append(paramsAnd, param)
			}
		}

		// Process OR filters
		for _, filter := range req.Filters.Or {
			whereClause, param := pagination.BuildWhereClause(filter)
			if whereClause != "" {
				whereOr = append(whereOr, whereClause)
				paramsOr = append(paramsOr, param)
			}
		}
	}

	// Combine all conditions
	var finalWhere string
	var finalParams []interface{}

	// Add AND conditions
	if len(whereAnd) > 0 {
		finalWhere = strings.Join(whereAnd, " AND ")
		finalParams = append(finalParams, paramsAnd...)
	}

	// Add OR conditions (wrapped in parentheses)
	if len(whereOr) > 0 {
		orClause := "(" + strings.Join(whereOr, " OR ") + ")"
		if finalWhere != "" {
			finalWhere += " AND " + orClause
		} else {
			finalWhere = orClause
		}
		finalParams = append(finalParams, paramsOr...)
	}

	// Apply the combined where clause
	if finalWhere != "" {
		query = query.Where(finalWhere, finalParams...)
	}

	// Get filtered count
	if err := query.Count(&total64).Error; err != nil {
		return nil, 0, 0, err
	}
	totalFiltered = int(total64)

	// Apply pagination and select
	err = query.Select(`
		id,
		_log_upload_id,
		idProyek,
		uraianJenisProyek,
		nib,
		tglDownload as tgl_download,
		tglDownloadExcel,
		tglTerbitOss,
		tglTerbitOssExcel,
		tglPengajuan,
		tglPengajuanExcel,
		lastUpdateProyek,
		lastUpdateProyekRaw,
		pendaftarNIK,
		pendaftarTglLahir,
		pendaftarGender,
		pendaftarNama,
		pendaftarTelp,
		pendaftarEmail,
		perusahaanNPWP,
		perusahaanNama,
		perusahaanAlamat,
		perusahaanKelurahan,
		perusahaanKecamatan,
		perusahaanKota,
		perusahaanProv,
		perusahaanLon,
		perusahaanLat,
		perusahaanSkala,
		perusahaanSkalaKbli,
		jenisBadan,
		jenisBadanDetail,
		statusNIB,
		statusPM,
		resiko,
		kbli,
		kbliJudul,
		sektorPembina,
		tenagaKerja,
		namaProyek,
		luasTanah,
		satuanTanah,
		invModalTetap,
		invMesinPeralatanImpor,
		invMesinPeralatan,
		invBeliPematanganTanah,
		invBangunanGedung,
		invModalKerja,
		invLain,
		invJumlah,
		invJumlahRumus,
		_created_at,
		_created_by,
		_updated_at,
		_updated_by,
		_input_manual
	`).Limit(req.PerPage).Offset((req.Page - 1) * req.PerPage).Scan(&data).Error

	if err != nil {
		return nil, 0, 0, err
	}

	return data, total, totalFiltered, nil
}
