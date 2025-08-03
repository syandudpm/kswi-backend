package oss

import (
	"kswi-backend/internal/shared/pagination"
	"time"
)

type DtDatabaseResponse struct {
	ID                     int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	LogUploadID            *int       `json:"log_upload_id" gorm:"column:_log_upload_id"`
	IdProyek               *string    `json:"id_proyek" gorm:"column:idProyek;size:150"`
	UraianJenisProyek      *string    `json:"uraian_jenis_proyek" gorm:"column:uraianJenisProyek;size:150"`
	NIB                    *string    `json:"nib" gorm:"column:nib;size:25"`
	TglDownload            *time.Time `json:"tgl_download" gorm:"column:tglDownload"`
	TglDownloadExcel       *string    `json:"tgl_download_excel" gorm:"column:tglDownloadExcel;size:10"`
	TglTerbitOss           *time.Time `json:"tgl_terbit_oss" gorm:"column:tglTerbitOss"`
	TglTerbitOssExcel      *string    `json:"tgl_terbit_oss_excel" gorm:"column:tglTerbitOssExcel;size:10"`
	TglPengajuan           *time.Time `json:"tgl_pengajuan" gorm:"column:tglPengajuan"`
	TglPengajuanExcel      *string    `json:"tgl_pengajuan_excel" gorm:"column:tglPengajuanExcel;size:10"`
	LastUpdateProyek       *time.Time `json:"last_update_proyek" gorm:"column:lastUpdateProyek"`
	LastUpdateProyekRaw    *string    `json:"last_update_proyek_raw" gorm:"column:lastUpdateProyekRaw;size:150"`
	PendaftarNIK           *string    `json:"pendaftar_nik" gorm:"column:pendaftarNIK;size:25"`
	PendaftarTglLahir      *time.Time `json:"pendaftar_tgl_lahir" gorm:"column:pendaftarTglLahir"`
	PendaftarGender        *string    `json:"pendaftar_gender" gorm:"column:pendaftarGender;size:25"`
	PendaftarNama          *string    `json:"pendaftar_nama" gorm:"column:pendaftarNama;size:245"`
	PendaftarTelp          *string    `json:"pendaftar_telp" gorm:"column:pendaftarTelp;size:445"`
	PendaftarEmail         *string    `json:"pendaftar_email" gorm:"column:pendaftarEmail;size:145"`
	PerusahaanNPWP         *string    `json:"perusahaan_npwp" gorm:"column:perusahaanNPWP;size:445"`
	PerusahaanNama         *string    `json:"perusahaan_nama" gorm:"column:perusahaanNama;size:545"`
	PerusahaanAlamat       *string    `json:"perusahaan_alamat" gorm:"column:perusahaanAlamat;size:545"`
	PerusahaanKelurahan    *string    `json:"perusahaan_kelurahan" gorm:"column:perusahaanKelurahan;size:445"`
	PerusahaanKecamatan    *string    `json:"perusahaan_kecamatan" gorm:"column:perusahaanKecamatan;size:445"`
	PerusahaanKota         *string    `json:"perusahaan_kota" gorm:"column:perusahaanKota;size:445"`
	PerusahaanProv         *string    `json:"perusahaan_prov" gorm:"column:perusahaanProv;size:445"`
	PerusahaanLon          *string    `json:"perusahaan_lon" gorm:"column:perusahaanLon;size:145"`
	PerusahaanLat          *string    `json:"perusahaan_lat" gorm:"column:perusahaanLat;size:145"`
	PerusahaanSkala        *string    `json:"perusahaan_skala" gorm:"column:perusahaanSkala;size:445"`
	PerusahaanSkalaKbli    *string    `json:"perusahaan_skala_kbli" gorm:"column:perusahaanSkalaKbli;size:445"`
	JenisBadan             *string    `json:"jenis_badan" gorm:"column:jenisBadan;size:445"`
	JenisBadanDetail       *string    `json:"jenis_badan_detail" gorm:"column:jenisBadanDetail;size:445"`
	StatusNIB              *string    `json:"status_nib" gorm:"column:statusNIB;size:445"`
	StatusPM               *string    `json:"status_pm" gorm:"column:statusPM;size:445"`
	Resiko                 *string    `json:"resiko" gorm:"column:resiko;size:445"`
	Kbli                   *string    `json:"kbli" gorm:"column:kbli;size:445"`
	KbliJudul              *string    `json:"kbli_judul" gorm:"column:kbliJudul;size:445"`
	SektorPembina          *string    `json:"sektor_pembina" gorm:"column:sektorPembina;size:445"`
	TenagaKerja            *int       `json:"tenaga_kerja" gorm:"column:tenagaKerja"`
	NamaProyek             *string    `json:"nama_proyek" gorm:"column:namaProyek;size:550"`
	LuasTanah              *string    `json:"luas_tanah" gorm:"column:luasTanah;size:20"`
	SatuanTanah            *string    `json:"satuan_tanah" gorm:"column:satuanTanah;size:20"`
	InvModalTetap          *uint64    `json:"inv_modal_tetap" gorm:"column:invModalTetap"`
	InvMesinPeralatanImpor *uint64    `json:"inv_mesin_peralatan_impor" gorm:"column:invMesinPeralatanImpor"`
	InvMesinPeralatan      *uint64    `json:"inv_mesin_peralatan" gorm:"column:invMesinPeralatan"`
	InvBeliPematanganTanah *uint64    `json:"inv_beli_pematangan_tanah" gorm:"column:invBeliPematanganTanah"`
	InvBangunanGedung      *uint64    `json:"inv_bangunan_gedung" gorm:"column:invBangunanGedung"`
	InvModalKerja          *uint64    `json:"inv_modal_kerja" gorm:"column:invModalKerja"`
	InvLain                *uint64    `json:"inv_lain" gorm:"column:invLain"`
	InvJumlah              *uint64    `json:"inv_jumlah" gorm:"column:invJumlah"`
	InvJumlahRumus         *uint64    `json:"inv_jumlah_rumus" gorm:"column:invJumlahRumus"`
	CreatedAt              *time.Time `json:"created_at" gorm:"column:_created_at;autoCreateTime"`
	CreatedBy              *int       `json:"created_by" gorm:"column:_created_by"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"column:_updated_at;autoUpdateTime"`
	UpdatedBy              *int       `json:"updated_by" gorm:"column:_updated_by"`
	InputManual            *int       `json:"input_manual" gorm:"column:_input_manual;default:0"`
}

type DtDatabaseRequest struct {
	pagination.PaginationRequest
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
}
