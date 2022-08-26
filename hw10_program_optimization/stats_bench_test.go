package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	reader, _ := zip.OpenReader("testdata/users.dat.zip")
	defer reader.Close()
	data, _ := reader.File[0].Open()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "com")
	}

	b.StopTimer()
	b.ReportAllocs()
}
