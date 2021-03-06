package artifactory

import (
	"context"
	"github.com/atlassian/go-artifactory/pkg/artifactory"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

func resourceArtifactoryLocalRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocalRepositoryCreate,
		Read:   resourceLocalRepositoryRead,
		Update: resourceLocalRepositoryUpdate,
		Delete: resourceLocalRepositoryDelete,
		Exists: resourceLocalRepositoryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "generic",
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"includes_pattern": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "**/*",
			},
			"excludes_pattern": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"repo_layout_ref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"handle_releases": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"handle_snapshots": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"max_unique_snapshots": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"debian_trivial_layout": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"checksum_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "client-checksums",
			},
			"max_unique_tags": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"snapshot_version_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "non-unique",
			},
			"suppress_pom_consistency_checks": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"blacked_out": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"property_sets": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"archive_browsing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"calculate_yum_metadata": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"yum_root_depth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"docker_api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "V2",
			},
			"enable_file_lists_indexing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func unmarshalLocalRepository(s *schema.ResourceData) *artifactory.LocalRepository {
	d := &ResourceData{s}

	repo := new(artifactory.LocalRepository)

	repo.RClass = artifactory.String("local")

	repo.Key = d.GetStringRef("key")
	repo.PackageType = d.GetStringRef("package_type")
	repo.Description = d.GetStringRef("description")
	repo.Notes = d.GetStringRef("notes")
	repo.DebianTrivialLayout = d.GetBoolRef("debian_trivial_layout")
	repo.IncludesPattern = d.GetStringRef("includes_pattern")
	repo.ExcludesPattern = d.GetStringRef("excludes_pattern")
	repo.RepoLayoutRef = d.GetStringRef("repo_layout_ref")
	repo.MaxUniqueTags = d.GetIntRef("max_unique_tags")
	repo.BlackedOut = d.GetBoolRef("blacked_out")
	repo.CalculateYumMetadata = d.GetBoolRef("calculate_yum_metadata")
	repo.YumRootDepth = d.GetIntRef("yum_root_depth")
	repo.ArchiveBrowsingEnabled = d.GetBoolRef("archive_browsing_enabled")
	repo.DockerApiVersion = d.GetStringRef("docker_api_verision")
	repo.EnableFileListsIndexing = d.GetBoolRef("enable_file_lists_indexing")
	repo.PropertySets = d.GetSetRef("property_sets")
	repo.HandleReleases = d.GetBoolRef("handle_releases")
	repo.HandleSnapshots = d.GetBoolRef("handle_snapshots")
	repo.ChecksumPolicyType = d.GetStringRef("checksum_policy_type")
	repo.MaxUniqueSnapshots = d.GetIntRef("max_unique_snapshots")
	repo.SnapshotVersionBehavior = d.GetStringRef("snapshot_version_behavior")
	repo.SuppressPomConsistencyChecks = d.GetBoolRef("suppress_pom_consistency_checks")

	return repo
}

func resourceLocalRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Client)

	repo := unmarshalLocalRepository(d)

	_, err := c.Repositories.CreateLocal(context.Background(), repo)
	if err != nil {
		return err
	}

	d.SetId(*repo.Key)
	return resourceLocalRepositoryRead(d, m)
}

func resourceLocalRepositoryRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Client)

	repo, resp, err := c.Repositories.GetLocal(context.Background(), d.Id())

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
	} else if err == nil {
		d.Set("key", repo.Key)
		d.Set("type", repo.RClass)
		d.Set("package_type", repo.PackageType)
		d.Set("description", repo.Description)
		d.Set("notes", repo.Notes)
		d.Set("includes_pattern", repo.IncludesPattern)
		d.Set("excludes_pattern", repo.ExcludesPattern)
		d.Set("repo_layout_ref", repo.RepoLayoutRef)
		d.Set("debian_trivial_layout", repo.DebianTrivialLayout)
		d.Set("max_unique_tags", repo.MaxUniqueTags)
		d.Set("blacked_out", repo.BlackedOut)
		d.Set("archive_browsing_enabled", repo.ArchiveBrowsingEnabled)
		d.Set("calculate_yum_metadata", repo.CalculateYumMetadata)
		d.Set("yum_root_depth", repo.YumRootDepth)
		d.Set("docker_api_version", repo.DockerApiVersion)
		d.Set("enable_file_lists_indexing", repo.EnableFileListsIndexing)
		d.Set("property_sets", schema.NewSet(schema.HashString, CastToInterfaceArr(*repo.PropertySets)))
		d.Set("handle_releases", repo.HandleReleases)
		d.Set("handle_snapshots", repo.HandleSnapshots)
		d.Set("checksum_policy_type", repo.ChecksumPolicyType)
		d.Set("max_unique_snapshots", repo.MaxUniqueSnapshots)
		d.Set("snapshot_version_behavior", repo.SnapshotVersionBehavior)
		d.Set("suppress_pom_consistency_checks", repo.SuppressPomConsistencyChecks)
	}

	return err
}

func resourceLocalRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Client)

	repo := unmarshalLocalRepository(d)
	_, err := c.Repositories.UpdateLocal(context.Background(), d.Id(), repo)

	if err != nil {
		return err
	}

	d.SetId(*repo.Key)
	return resourceLocalRepositoryRead(d, m)
}

func resourceLocalRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Client)
	repo := unmarshalLocalRepository(d)

	resp, err := c.Repositories.DeleteLocal(context.Background(), *repo.Key)

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return err
}

func resourceLocalRepositoryExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*artifactory.Client)

	_, resp, err := c.Repositories.GetLocal(context.Background(), d.Id())

	// Cannot check for 404 because artifactory returns 400
	if resp.StatusCode == http.StatusBadRequest {
		return false, nil
	}

	return true, err
}
