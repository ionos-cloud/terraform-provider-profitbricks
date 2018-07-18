package profitbricks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
	"golang.org/x/crypto/ssh"
)

func resourceProfitBricksServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksServerCreate,
		Read:   resourceProfitBricksServerRead,
		Update: resourceProfitBricksServerUpdate,
		Delete: resourceProfitBricksServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProfitBricksServerImport,
		},
		Schema: map[string]*schema.Schema{

			//Server parameters
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"licence_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"boot_volume": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"boot_cdrom": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_family": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"boot_image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"primary_nic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewallrule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_key_path": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"image_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"licence_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bus": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"nic": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lan": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if new == "" {
									return true
								}
								return false
							},
						},
						"ips": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"nat": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"firewall_active": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"firewall": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"protocol": {
										Type:     schema.TypeString,
										Required: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											if strings.ToLower(old) == strings.ToLower(new) {
												return true
											}
											return false
										},
									},
									"source_mac": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ips": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"port_range_start": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 && v.(int) > 65534 {
												errors = append(errors, fmt.Errorf("Port start range must be between 1 and 65534"))
											}
											return
										},
									},

									"port_range_end": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 && v.(int) > 65534 {
												errors = append(errors, fmt.Errorf("Port end range must be between 1 and 65534"))
											}
											return
										},
									},
									"icmp_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"icmp_code": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	var image_alias string
	request := profitbricks.Server{
		Properties: profitbricks.ServerProperties{
			Name:  d.Get("name").(string),
			Cores: d.Get("cores").(int),
			RAM:   d.Get("ram").(int),
		},
	}
	dcId := d.Get("datacenter_id").(string)

	isSnapshot := false
	if v, ok := d.GetOk("availability_zone"); ok {
		request.Properties.AvailabilityZone = v.(string)
	}

	if v, ok := d.GetOk("cpu_family"); ok {
		if v.(string) != "" {
			request.Properties.CPUFamily = v.(string)
		}
	}
	if vRaw, ok := d.GetOk("volume"); ok {

		volumeRaw := vRaw.(*schema.Set).List()

		for _, raw := range volumeRaw {
			rawMap := raw.(map[string]interface{})
			var imagePassword string
			//Can be one file or a list of files
			var sshkey_path []interface{}
			var image, licenceType, availabilityZone string

			imagePassword = d.Get("image_password").(string)
			sshkey_path = d.Get("ssh_key_path").([]interface{})

			image_name := d.Get("boot_image").(string)
			if !IsValidUUID(image_name) {
				img, err := getImage(client, dcId, image_name, rawMap["disk_type"].(string))
				if err != nil {
					return err
				}

				if img != nil {
					image = img.ID
				}

				//if no image id was found with that name we look for a matching snapshot
				if image == "" {
					image = getSnapshotId(client, image_name)
					if image != "" {
						isSnapshot = true
					} else {
						dc, err := client.GetDatacenter(dcId)
						if err != nil {
							return fmt.Errorf("Error fetching datacenter %s: (%s)", dcId, err)
						}
						image_alias = getImageAlias(client, image_name, dc.Properties.Location)
					}
				}
				if image == "" && image_alias == "" {
					return fmt.Errorf("Could not find an image/imagealias/snapshot that matches %s ", image_name)
				}
				if imagePassword == "" && len(sshkey_path) == 0 && isSnapshot == false && img.Properties.Public {
					return fmt.Errorf("Either 'image_password' or 'ssh_key_path' must be provided.")
				}
			} else {
				img, err := client.GetImage(image_name)

				var apiError profitbricks.ApiError

				if apiError, ok = err.(profitbricks.ApiError); !ok {
					return fmt.Errorf("Error fetching image %s: (%s)", image_name, err)
				}

				if apiError.HttpStatusCode() == 404 {
					img, err := client.GetSnapshot(image_name)

					if apiError, ok := err.(profitbricks.ApiError); !ok {
						if apiError.HttpStatusCode() == 404 {
							return fmt.Errorf("image/snapshot: %s Not Found", img.Response)
						}
					}

					isSnapshot = true
				} else {
					if err != nil {
						return fmt.Errorf("Error fetching image/snapshot: %s", err)
					}
				}
				if img.Properties.Public == true && isSnapshot == false {
					if imagePassword == "" && len(sshkey_path) == 0 {
						return fmt.Errorf("Either 'image_password' or 'ssh_key_path' must be provided.")
					}
					img, err := getImage(client, d.Get("datacenter_id").(string), image_name, rawMap["disk_type"].(string))
					if err != nil {
						return err
					}
					if img != nil {
						image = img.ID
					}
				} else {
					img, err := client.GetImage(image_name)
					if err != nil {
						img, err := client.GetSnapshot(image_name)
						if err != nil {
							return fmt.Errorf("Error fetching image/snapshot: %s", img.Response)
						}
						isSnapshot = true
					}
					if img.Properties.Public == true && isSnapshot == false {
						if imagePassword == "" && len(sshkey_path) == 0 {
							return fmt.Errorf("Either 'image_password' or 'ssh_key_path' must be provided.")
						}
						image = image_name
					} else {
						image = image_name
					}
				}
			}

			if rawMap["licence_type"] != nil {
				licenceType = rawMap["licence_type"].(string)
			}

			var publicKeys []string
			if len(sshkey_path) != 0 {
				for _, path := range sshkey_path {
					log.Printf("[DEBUG] Reading file %s", path)
					publicKey, err := readPublicKey(path.(string))
					if err != nil {
						return fmt.Errorf("Error fetching sshkey from file (%s) %s", path, err.Error())
					}
					publicKeys = append(publicKeys, publicKey)
				}
			}
			if rawMap["availability_zone"] != nil {
				availabilityZone = rawMap["availability_zone"].(string)
			}
			if image == "" && licenceType == "" && image_alias == "" && !isSnapshot {
				return fmt.Errorf("Either 'image', 'licenceType', or 'imageAlias' must be set.")
			}

			if isSnapshot == true && (imagePassword != "" || len(publicKeys) > 0) {
				return fmt.Errorf("Passwords/SSH keys  are not supported for snapshots.")
			}

			request.Entities = &profitbricks.ServerEntities{
				Volumes: &profitbricks.Volumes{
					Items: []profitbricks.Volume{
						{
							Properties: profitbricks.VolumeProperties{
								Name:             rawMap["name"].(string),
								Size:             rawMap["size"].(int),
								Type:             rawMap["disk_type"].(string),
								ImagePassword:    imagePassword,
								Image:            image,
								ImageAlias:       image_alias,
								Bus:              rawMap["bus"].(string),
								LicenceType:      licenceType,
								AvailabilityZone: availabilityZone,
							},
						},
					},
				},
			}

			if len(publicKeys) == 0 {
				request.Entities.Volumes.Items[0].Properties.SSHKeys = nil
			} else {
				request.Entities.Volumes.Items[0].Properties.SSHKeys = publicKeys
			}
		}

	}

	if _, ok := d.GetOk("nic"); ok {
		nic := profitbricks.Nic{Properties: &profitbricks.NicProperties{
			Lan: d.Get("nic.0.lan").(int),
		}}

		if v, ok := d.GetOk("nic.0.name"); ok {
			nic.Properties.Name = v.(string)
		}

		if v, ok := d.GetOk("nic.0.dhcp"); ok {
			val := v.(bool)
			nic.Properties.Dhcp = &val
		}

		if v, ok := d.GetOk("nic.0.firewall_active"); ok {
			nic.Properties.FirewallActive = v.(bool)
		}

		if v, ok := d.GetOk("nic.0.ip"); ok {
			ips := strings.Split(v.(string), ",")
			if len(ips) > 0 {
				nic.Properties.Ips = ips
			}
		}

		if v, ok := d.GetOk("nic.0.nat"); ok {
			nic.Properties.Nat = v.(bool)
		}

		request.Entities.Nics = &profitbricks.Nics{
			Items: []profitbricks.Nic{
				nic,
			},
		}

		if _, ok := d.GetOk("nic.0.firewall"); ok {
			firewall := profitbricks.FirewallRule{
				Properties: profitbricks.FirewallruleProperties{
					Protocol: d.Get("nic.0.firewall.0.protocol").(string),
				},
			}

			if v, ok := d.GetOk("nic.0.firewall.0.name"); ok {
				firewall.Properties.Name = v.(string)
			}

			if v, ok := d.GetOk("nic.0.firewall.0.source_mac"); ok {
				val := v.(string)
				firewall.Properties.SourceMac = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.source_ip"); ok {
				val := v.(string)
				firewall.Properties.SourceIP = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.target_ip"); ok {
				val := v.(string)
				firewall.Properties.TargetIP = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.port_range_start"); ok {
				val := v.(int)
				firewall.Properties.PortRangeStart = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.port_range_end"); ok {
				val := v.(int)
				firewall.Properties.PortRangeEnd = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.icmp_type"); ok {
				tempIcmpType := v.(string)
				if tempIcmpType != "" {
					i, _ := strconv.Atoi(tempIcmpType)
					firewall.Properties.IcmpType = &i
				}
			}
			if v, ok := d.GetOk("nic.0.firewall.0.icmp_code"); ok {
				tempIcmpCode := v.(string)
				if tempIcmpCode != "" {
					i, _ := strconv.Atoi(tempIcmpCode)
					firewall.Properties.IcmpCode = &i
				}
			}
			request.Entities.Nics.Items[0].Entities = &profitbricks.NicEntities{
				FirewallRules: &profitbricks.FirewallRules{
					Items: []profitbricks.FirewallRule{
						firewall,
					},
				},
			}
		}
	}

	if len(request.Entities.Nics.Items[0].Properties.Ips) == 0 {
		request.Entities.Nics.Items[0].Properties.Ips = nil
	}
	server, err := client.CreateServer(d.Get("datacenter_id").(string), request)

	jsn, _ := json.Marshal(request)
	log.Println("[DEBUG] Server request", string(jsn))
	log.Println("[DEBUG] Server response", server.Response)

	if err != nil {
		return fmt.Errorf(
			"Error creating server: (%s)", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, server.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId(server.ID)
	server, err = client.GetServer(d.Get("datacenter_id").(string), server.ID)
	if err != nil {
		return fmt.Errorf("Error fetching server: (%s)", err)
	}

	firewallRules, err := client.ListFirewallRules(d.Get("datacenter_id").(string), server.ID, server.Entities.Nics.Items[0].ID)

	if len(firewallRules.Items) > 0 {
		d.Set("firewallrule_id", firewallRules.Items[0].ID)
	}

	d.Set("primary_nic", server.Entities.Nics.Items[0].ID)
	if len(server.Entities.Nics.Items[0].Properties.Ips) > 0 {
		d.SetConnInfo(map[string]string{
			"type":     "ssh",
			"host":     server.Entities.Nics.Items[0].Properties.Ips[0],
			"password": request.Entities.Volumes.Items[0].Properties.ImagePassword,
		})
	}
	return resourceProfitBricksServerRead(d, meta)
}

func resourceProfitBricksServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	server, err := client.GetServer(dcId, serverId)
	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a server ID %s %s", d.Id(), err)
	}
	d.Set("name", server.Properties.Name)
	d.Set("cores", server.Properties.Cores)
	d.Set("ram", server.Properties.RAM)
	d.Set("availability_zone", server.Properties.AvailabilityZone)
	d.Set("cpu_family", server.Properties.CPUFamily)

	if primarynic, ok := d.GetOk("primary_nic"); ok {
		d.Set("primary_nic", primarynic.(string))

		nic, err := client.GetNic(dcId, serverId, primarynic.(string))
		if err != nil {
			return fmt.Errorf("Error occured while fetching nic %s for server ID %s %s", primarynic.(string), d.Id(), err)
		}

		if len(nic.Properties.Ips) > 0 {
			d.Set("primary_ip", nic.Properties.Ips[0])
		}

		network := map[string]interface{}{
			"lan":             nic.Properties.Lan,
			"name":            nic.Properties.Name,
			"dhcp":            *nic.Properties.Dhcp,
			"nat":             nic.Properties.Nat,
			"firewall_active": nic.Properties.FirewallActive,
			"ips":             nic.Properties.Ips,
		}

		if len(nic.Properties.Ips) > 0 {
			network["ip"] = nic.Properties.Ips[0]
		}

		if firewall_id, ok := d.GetOk("firewallrule_id"); ok {
			firewall, err := client.GetFirewallRule(dcId, serverId, primarynic.(string), firewall_id.(string))
			if err != nil {
				return fmt.Errorf("Error occured while fetching firewallrule %s for server ID %s %s", firewall_id.(string), serverId, err)
			}

			fw := map[string]interface{}{
				"protocol": firewall.Properties.Protocol,
				"name":     firewall.Properties.Name,
			}

			if firewall.Properties.SourceMac != nil {
				fw["source_mac"] = *firewall.Properties.SourceMac
			}

			if firewall.Properties.SourceIP != nil {
				fw["source_ip"] = *firewall.Properties.SourceIP
			}

			if firewall.Properties.TargetIP != nil {
				fw["target_ip"] = *firewall.Properties.TargetIP
			}

			if firewall.Properties.PortRangeStart != nil {
				fw["port_range_start"] = *firewall.Properties.PortRangeStart
			}

			if firewall.Properties.PortRangeEnd != nil {
				fw["port_range_end"] = *firewall.Properties.PortRangeEnd
			}

			if firewall.Properties.IcmpType != nil {
				fw["icmp_type"] = *firewall.Properties.IcmpType
			}

			if firewall.Properties.IcmpCode != nil {
				fw["icmp_code"] = *firewall.Properties.IcmpCode
			}

			network["firewall"] = []map[string]interface{}{fw}
		}

		networks := []map[string]interface{}{network}
		if err := d.Set("nic", networks); err != nil {
			return fmt.Errorf("[ERROR] unable saving nic to state ProfitBricks Server (%s): %s", serverId, err)
		}
	}

	if server.Properties.BootVolume != nil {
		d.Set("boot_volume", server.Properties.BootVolume.ID)

		volumeObj, err := client.GetAttachedVolume(dcId, serverId, server.Properties.BootVolume.ID)
		if err != nil {
			return fmt.Errorf("Error occured while fetching attached volume %s from server ID %s %s", server.Properties.BootVolume.ID, serverId, err)
		}

		volumeItem := map[string]interface{}{
			"name":      volumeObj.Properties.Name,
			"disk_type": volumeObj.Properties.Type,
			"size":      volumeObj.Properties.Size,
		}

		volumesList := []map[string]interface{}{volumeItem}
		if err := d.Set("volume", volumesList); err != nil {
			return fmt.Errorf("[DEBUG] Error saving volume to state for ProfitBricks server (%s): %s", d.Id(), err)
		}
	}

	if server.Properties.BootCdrom != nil {
		d.Set("boot_cdrom", server.Properties.BootCdrom.ID)
	}
	return nil
}

func resourceProfitBricksServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)

	request := profitbricks.ServerProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		request.Name = n.(string)
	}
	if d.HasChange("cores") {
		_, n := d.GetChange("cores")
		request.Cores = n.(int)
	}
	if d.HasChange("ram") {
		_, n := d.GetChange("ram")
		request.RAM = n.(int)
	}
	if d.HasChange("availability_zone") {
		_, n := d.GetChange("availability_zone")
		request.AvailabilityZone = n.(string)
	}
	if d.HasChange("cpu_family") {
		_, n := d.GetChange("cpu_family")
		request.CPUFamily = n.(string)
	}
	server, err := client.UpdateServer(dcId, d.Id(), request)

	if err != nil {
		return fmt.Errorf("Error occured while updating server ID %s %s", d.Id(), err)
	}

	_, errState := getStateChangeConf(meta, d, server.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}
	//Volume stuff
	if d.HasChange("volume") {
		_, new := d.GetChange("volume")

		newVolume := new.(*schema.Set).List()
		properties := profitbricks.VolumeProperties{}

		for _, raw := range newVolume {
			rawMap := raw.(map[string]interface{})
			if rawMap["name"] != nil {
				properties.Name = rawMap["name"].(string)
			}
			if rawMap["size"] != nil {
				properties.Size = rawMap["size"].(int)
			}
			if rawMap["bus"] != nil {
				properties.Bus = rawMap["bus"].(string)
			}
		}

		volume, err := client.UpdateVolume(d.Get("datacenter_id").(string), server.Entities.Volumes.Items[0].ID, properties)

		if err != nil {
			return fmt.Errorf("Error patching volume (%s) (%s)", d.Id(), err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, volume.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}
	}

	//Nic stuff
	if d.HasChange("nic") {
		nic := &profitbricks.Nic{}
		for _, n := range server.Entities.Nics.Items {
			if n.ID == d.Get("primary_nic").(string) {
				nic = &n
				break
			}
		}

		properties := profitbricks.NicProperties{
			Lan: d.Get("nic.0.lan").(int),
		}

		if v, ok := d.GetOk("nic.0.name"); ok {
			properties.Name = v.(string)
		}

		if v, ok := d.GetOk("nic.0.ip"); ok {
			ips := strings.Split(v.(string), ",")
			if len(ips) > 0 {
				nic.Properties.Ips = ips
			}
		}

		if v, ok := d.GetOk("nic.0.dhcp"); ok {
			val := v.(bool)
			properties.Dhcp = &val
		}

		if v, ok := d.GetOk("nic.0.nat"); ok {
			properties.Nat = v.(bool)
		}

		nic, err := client.UpdateNic(d.Get("datacenter_id").(string), server.ID, nic.ID, properties)

		if err != nil {
			return fmt.Errorf(
				"Error updating nic (%s)", err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, nic.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}
	}

	return resourceProfitBricksServerRead(d, meta)
}

func resourceProfitBricksServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)

	server, err := client.GetServer(dcId, d.Id())

	if err != nil {
		return fmt.Errorf("Error occured while fetching a server ID %s %s", d.Id(), err)
	}

	if server.Properties.BootVolume != nil {
		resp, err := client.DeleteVolume(dcId, server.Properties.BootVolume.ID)
		if err != nil {
			return fmt.Errorf("Error occured while delete volume %s of server ID %s %s", server.Properties.BootVolume.ID, d.Id(), err)
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
		if errState != nil {
			return errState
		}
	}

	resp, err := client.DeleteServer(dcId, d.Id())
	if err != nil {
		return fmt.Errorf("An error occured while deleting a server ID %s %s", d.Id(), err)

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}

//Reads public key from file and returns key string iff valid
func readPublicKey(path string) (key string, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return "", err
	}
	return string(ssh.MarshalAuthorizedKey(pubKey)[:]), nil
}
