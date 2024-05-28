// Copyright 2024 Stefan Prodan.
// SPDX-License-Identifier: AGPL-3.0

package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fluxcd/pkg/apis/kustomize"
	. "github.com/onsi/gomega"
	cp "github.com/otiai10/copy"
	"sigs.k8s.io/yaml"
)

func TestBuild_Defaults(t *testing.T) {
	g := NewWithT(t)
	const version = "v1.3.0"
	options := MakeDefaultOptions()
	options.Version = version

	srcDir := filepath.Join("testdata", version)
	goldenFile := filepath.Join("testdata", version+"-golden", "default.kustomization.yaml")

	dstDir, err := testTempDir(t)
	g.Expect(err).NotTo(HaveOccurred())

	objects, err := Build(srcDir, dstDir, options)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(objects).NotTo(BeEmpty())

	if os.Getenv("GEN_GOLDEN") == "true" {
		err = cp.Copy(filepath.Join(dstDir, "kustomization.yaml"), goldenFile)
		g.Expect(err).NotTo(HaveOccurred())
	}

	genK, err := os.ReadFile(filepath.Join(dstDir, "kustomization.yaml"))
	g.Expect(err).NotTo(HaveOccurred())

	goldenK, err := os.ReadFile(goldenFile)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(string(genK)).To(Equal(string(goldenK)))
}

func TestBuild_Patches(t *testing.T) {
	g := NewWithT(t)
	const version = "v1.3.0"
	options := MakeDefaultOptions()
	options.Version = version

	srcDir := filepath.Join("testdata", version)
	goldenFile := filepath.Join("testdata", version+"-golden", "patches.kustomization.yaml")

	dstDir, err := testTempDir(t)
	g.Expect(err).NotTo(HaveOccurred())

	patches := []kustomize.Patch{
		{
			Patch: `
- op: remove
  path: /metadata/labels/pod-security.kubernetes.io~1warn
- op: remove
  path: /metadata/labels/pod-security.kubernetes.io~1warn-version
`,
			Target: &kustomize.Selector{
				Kind: "Namespace",
			},
		},
	}
	patchesData, err := yaml.Marshal(patches)
	g.Expect(err).NotTo(HaveOccurred())
	options.Patches = string(patchesData)

	objects, err := Build(srcDir, dstDir, options)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(objects).NotTo(BeEmpty())

	if os.Getenv("GEN_GOLDEN") == "true" {
		err = cp.Copy(filepath.Join(dstDir, "kustomization.yaml"), goldenFile)
		g.Expect(err).NotTo(HaveOccurred())
	}

	genK, err := os.ReadFile(filepath.Join(dstDir, "kustomization.yaml"))
	g.Expect(err).NotTo(HaveOccurred())

	goldenK, err := os.ReadFile(goldenFile)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(string(genK)).To(Equal(string(goldenK)))

	found := false
	for _, obj := range objects {
		if obj.GetKind() == "Namespace" {
			found = true
			labels := obj.GetLabels()
			g.Expect(labels).NotTo(HaveKey("pod-security.kubernetes.io/warn"))
			g.Expect(labels).NotTo(HaveKey("pod-security.kubernetes.io/warn-version"))
			g.Expect(obj.GetAnnotations()).To(HaveKeyWithValue("fluxcd.controlplane.io/prune", "disabled"))
		}
	}
	g.Expect(found).To(BeTrue())
}

func TestBuild_InvalidPatches(t *testing.T) {
	g := NewWithT(t)
	const version = "v1.3.0"
	options := MakeDefaultOptions()
	options.Version = version
	srcDir := filepath.Join("testdata", version)

	dstDir, err := testTempDir(t)
	g.Expect(err).NotTo(HaveOccurred())

	patches := []kustomize.Patch{
		{
			Patch: `
- op: removes
  path: /metadata/labels/pod-security.kubernetes.io~1warn
`,
			Target: &kustomize.Selector{
				Kind: "Namespace",
			},
		},
	}
	patchesData, err := yaml.Marshal(patches)
	g.Expect(err).NotTo(HaveOccurred())
	options.Patches = string(patchesData)

	_, err = Build(srcDir, dstDir, options)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("Unexpected kind: removes"))
}

func testTempDir(t *testing.T) (string, error) {
	tmpDir := t.TempDir()

	tmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		return "", fmt.Errorf("error evaluating symlink: '%w'", err)
	}

	return tmpDir, err
}