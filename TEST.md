How to launch tests
===================

	go test ./...

Note
====

To use appenginetesting, edit `appengine_internal/api_dev.go`.

Line 57:

	if c != nil {
		instanceConfig.AppID = string(c.AppId)
		instanceConfig.APIPort = int(*c.ApiPort)
		instanceConfig.VersionID = string(c.VersionId)
		instanceConfig.InstanceID = *c.InstanceId
		instanceConfig.Datacenter = *c.Datacenter
	} else {
		instanceConfig.AppID = "testapp"
		instanceConfig.APIPort = 0
		instanceConfig.VersionID = "1.7.7"
		instanceConfig.InstanceID = "instanceid"
		instanceConfig.Datacenter = "instanceid"
	}
  
Line 132:

	if len(raw) == 0 {
		return nil
	}

Line 260:

	//func init() { os.DisableWritesForAppEngine = true }

