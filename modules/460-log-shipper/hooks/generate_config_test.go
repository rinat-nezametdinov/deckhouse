/*
Copyright 2021 Flant CJSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hooks

import (
	"encoding/base64"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Log shipper :: generate config from crd ::", func() {
	f := HookExecutionConfigInit(`{"logShipper": {"internal": {"activated": false}}}`, ``)
	f.RegisterCRD("deckhouse.io", "v1alpha1", "ClusterLoggingConfig", false)
	f.RegisterCRD("deckhouse.io", "v1alpha1", "ClusterLogDestination", false)
	f.RegisterCRD("deckhouse.io", "v1alpha1", "PodLoggingConfig", true)

	Context("Simple pair", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSetAndWaitForBindingContexts(`
apiVersion: deckhouse.io/v1alpha1
kind: ClusterLoggingConfig
metadata:
  name: test-source
spec:
  type: KubernetesPods
  kubernetesPods:
    namespaceSelector:
      matchNames:
      - tests-whispers
    labelSelector:
      matchLabels:
        app: test
  destinationRefs:
    - test-es-dest
    - test-loki-dest
    - test-logstash-dest
---
apiVersion: deckhouse.io/v1alpha1
kind: ClusterLogDestination
metadata:
  name: test-loki-dest
spec:
  type: Loki
  loki:
    endpoint: http://192.168.1.1:9000
  extraLabels:
    foo: bar
---
apiVersion: deckhouse.io/v1alpha1
kind: ClusterLogDestination
metadata:
  name: test-logstash-dest
spec:
  type: Logstash
  logstash:
    endpoint: 192.168.199.252:9009
    tls:
      caFile: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3ekNDQWFzQ0ZDalVzcGp5b29wVmdOcjR0TE5SS2hSWERmQXhNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1CNHgKQ3pBSkJnTlZCQVlUQWxKVk1ROHdEUVlEVlFRRERBWlVaWE4wUTBFd0hoY05NakV3TmpJeU1URTBOakEwV2hjTgpORGd4TVRBM01URTBOakEwV2pBZU1Rc3dDUVlEVlFRR0V3SlNWVEVQTUEwR0ExVUVBd3dHVkdWemRFTkJNSUlCCklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEzbG42U3pWSVR1VndlRFRneXR4TDZOTEMKditaeWc5d1dpVllSVnFjZ2hPU0FQMlhSZTJjTWJpYU5vbk9oZW00NDRka0JFY3d4WWhYZVhBWUE0N1dCSHZRRworWkZLOW9KaUJNZGRpSFpmNWpUV1pDK29KKzZMK0h0R2R4MUs3czNZaDM4aUMyWHRqelU5UUJzZmVCZUpIellZCmVXcm1MdDZpTjZRdDQ0Y3l3UHRKVW93ampKaU9YUHYxejluVDdjL3NGLzlTMUVsWENMV1B5dHdKV1NiMGVEUisKYTFGdmdFS1dxTWFySnJFbTFpWVhLU1FZUGFqWE9UU2hHaW9ITVZDK2VzMW55cHN6TG93ZUJ1Vjc5SS9WVnY0YQpnVk5CYTcwaWJEcXM3L3czcTJ3Q2I1ZlpBREU4MzJTcldIdGNtL0luSkNrQUt5czBySTlmODlQWHlHb1lNd0lECkFRQUJNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUM0b3lqL3V0VlFZa242eXU1UTBNbmVPK1YvTlNFSHhqTnIKcldOZnJuT2NTV2I4akFRWjN2ZFpHS1FMVW9raGFTUUNKQndhckxiQWlOVW1udG9ndEhEbEtnZUdxdGdVN3hWeQpJaTFCSlc1VkxyejhMNUdNREdQR2NmUjNpVDdKaDVyelM1UUc5YXlzVHgvMGpWaFN0T1I1cnFqdDlocmZrK0kvClQrT01QTTVrbHpzYXlnZTlkSEx1K3l1VzBzeHhHUk83KzlPeVY3bk9KNEd0TEhicWV0ajBWQUIraWpDMHp1NU0KakxDdm9aZEpQUFViWmVRenFlVW5ZTUwrQ0NERXpCSkdJRk9Xd2w1M2VTblFXbFdVaVJPZWNhd0hobkJzMWlHYgpTQ1BEMTFNMzRRRWZYMHBqQ054RUlzTUtvdFR6V2hFaCsvb0tyQnl2dW16SmpWeWtyU2l5Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
      crtFile: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN0akNDQVo0Q0ZHWDNFQ3I0V3dvVlBhUFpDNGZab042c2JYY09NQTBHQ1NxR1NJYjNEUUVCQ3dVQU1CNHgKQ3pBSkJnTlZCQVlUQWxKVk1ROHdEUVlEVlFRRERBWlVaWE4wUTBFd0hoY05NakV3TmpJeU1URTFOekUyV2hjTgpNelV3TXpBeE1URTFOekUyV2pBUk1ROHdEUVlEVlFRRERBWjJaV04wYjNJd2dnRWlNQTBHQ1NxR1NJYjNEUUVCCkFRVUFBNElCRHdBd2dnRUtBb0lCQVFER0JkSHBvWC9mQytaUkdFQVZpT2tyeE91b0JIazEyYVNLRldVU2hJSFcKZWowNC9zMUtjZFF5RUxlSlk5YUMxTzVuZ1hzdVpDVUNmS1NWdHE1Y3IySTV6cjRaaXNyM0JZK3JlcVBVYkVlYgpLNFBCdEVROUlibno2RTZMVUt3SitIRTFZamliRUFuRkRlamhSUWp6MHFUNWFYR1lNd0RkK1dGMUZ2YzFlUHkvCjhsZEc3YzNvRmczb0ZiV1p6bm9WQmYzOXh3WWZZdEZ2cGN2NWYwbW1SVmZlempRUk9nblhjT1dGb1F4VWcwSjEKV1FFM0xVSUdYMTBzQVpzdUpwMzVSN0tBL1pIRjZHcjhwemZIUmNRaHZPb2VBY0pPdTZZMFBaMnBwSzBhekt6LwpxeHMrZi9hUUJmc0N0c3V2Ty9HbmIvWWFDM1R3QTJmZXhlKzJBWjZGK1NBVEFnTUJBQUV3RFFZSktvWklodmNOCkFRRUxCUUFEZ2dFQkFFeEhkOUtBdkFZYTB2aG1aU0VkR1g3TnZIajhBWDFPV1VBcXZicHJ3YkZ1QkgyZm5LWCsKTmJGVHZXakpDUDdkem10cHphMVQ5RG1vOTJDNC9sWjk0Vy9Vc0pPRjJjSEFRUHlKdk5TdmJPVEg5YTAzajhCaAppbVJ3Zm0rTHNub3RGS3h3VTRhUCtRSEcrRVB2L0FDMDF3UDVhOWVpMEVZWnJIUXh1dTVsOWdURFdjU3Rra1o5Ci8xdzRFWGdNQ2xZVVdnQ1VHUTYvNy9XTkJONTNjWWZ5aU1QcS9VTmVQZUlhUkJDbXJxbklaUCtTWjVwMzFFUXMKZnIyak1rUUo5bTdqNlhWL0RrZFhTSWwrVmdmaVhRSXJDcVN2UXV3RldwdnBicFRPcFJOclhhNGlrMEJLMG1LaQpiYmkwTFVnbzJTcGJuSGlydGlWeVAvMTBCdWhmM3dISUdHUT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
      keyFile: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBeGdYUjZhRi8zd3ZtVVJoQUZZanBLOFRycUFSNU5kbWtpaFZsRW9TQjFubzlPUDdOClNuSFVNaEMzaVdQV2d0VHVaNEY3TG1RbEFueWtsYmF1WEs5aU9jNitHWXJLOXdXUHEzcWoxR3hIbXl1RHdiUkUKUFNHNTgraE9pMUNzQ2ZoeE5XSTRteEFKeFEzbzRVVUk4OUtrK1dseG1ETUEzZmxoZFJiM05Yajh2L0pYUnUzTgo2QllONkJXMW1jNTZGUVg5L2NjR0gyTFJiNlhMK1g5SnBrVlgzczQwRVRvSjEzRGxoYUVNVklOQ2RWa0JOeTFDCkJsOWRMQUdiTGlhZCtVZXlnUDJSeGVocS9LYzN4MFhFSWJ6cUhnSENUcnVtTkQyZHFhU3RHc3lzLzZzYlBuLzIKa0FYN0FyYkxyenZ4cDIvMkdndDA4QU5uM3NYdnRnR2VoZmtnRXdJREFRQUJBb0lCQURVcXd0MXpteDJMMkY3VgpuLzhvTDFLdElJaVFDdXRHY0VNUzAzeFJUM3NDZndXYWhBd0UyL0JGUk1JQ3FFbWdXaEk0VlpaekZPekNBbjZmCitkaXd6akt2SzZNMy9KNnVRNURLOE1uTCtMM1V4Ujl4QXhGV3lOS1FBT2F1MWtJbkRsNUM3T2ZWT29wSjNjajkKL0JWYTdTaDZBeUhXTDlscFo1MUVlVU5HSkxaMEpadWZCMVFiQVdpME5hRVpIdWFPL1FDWU55Qjh5Tk1PQkd5YQpPOUxtZHlDZk85VC9ZTFpXeC9kQ041WldZckhqVEpaREd3T3lCd1k1QjAzUWFmSitxQU5OSkVTTWV6bnlUdkRKCjk5d2hIQ0lxRjRDaHAwM2Y3Sm5QUXJCSDBIbWNDMW9BZjhMWFg5djEvdzY4Smpld1U3VUhoMzlWcTZ0NGNWZXAKdlh4YVdJRUNnWUVBN2dDTFNTVlJQUXFvRlBBcHhEMDVmQmpNUmd2M2tTbWlwWlVNOW5XMkR2WHNUUlFDVFNTcwpVL2JUMG5xZ0FtVTdXZVI3aUFMM2VKMU5ucjd5alc4ZUxaeXNGWUpvMzJNMmxHUGdIdVZoelJYL3ZuQ05CMUNHCmRrWVh5ZDVyK0grdkk1ZWxIcG8rbFVpYWd2NEtiQmtsQkNnRDllNFd6ZFhXN3F4STljc01PRU1DZ1lFQTFQOVIKeGhGNUJoNGVHV1g3RW1DMFRmMlVDa09wOTF1QXpQZDNmNFNQWHlkS2xxMDJCa3BCeFZKZEN2QVc2WlRGZ3FNdQp0Z1BxRi8rSzRNNy9IRStiODhoNytWdkJNVTIwdHFuNWM1Q2J0TUdlSU04MWkvdWxFODlqUlZ2LzI0Y3hZRitDCmlUdFZwUnh1NElNc05rdnAwNHhCMjZ1cGhHMk5HN0NVY2ZBdEkvRUNnWUVBcmpYQnZvbk5QRFFuc2lQVlBxcGUKQUlNYVN3K0phRDBrcTdVOVpzM2t0SEM0UmZjbWRCY3ErTTdNWDkyWWNBaHZlQzR4YWU1Wi9IU1FFMm5MbTFGQgpzcnRpanVBRktiYXloYzNSaUd2NHVhaW5xVnN6TDY1MnJlNUNqV1g4ZkVuaUJkaURhYklYcXlnWXlWZHdnNDJvCk5iR2dySXhaTHRPZTN0ZEhGSHRLOTRjQ2dZQnFXQ09xNGJSc0lvTmlxUEVuSnRNL0VUbGx1b3pVN0lHdFZHejgKWk9IMFh6aTFiRHZKL2k5Q1pySC9zUW12aTlEbFBiWW51R0tib3NIakpsWm0relJoRGhzZnovandOZHpoU3BJNgphZHZqN3J1Vm8vOFhLZ2dza09IK2trVjNoTk5aUzdadjhBajl5K2xyL1BJSkZmUGo1R1pKV0RibDRKQ1FYNlJ1CkVyMW04UUtCZ0VJdE5JSktDOEtNcjJ4VlBjbmo1NExZZ1BvYnhRcktLTlNnRUMrRTNkRFY4TEQyNnZHSmZRY0kKTDBsUE8zVm1vWWRaQnlraUF0NUNYRzUvRks5SkNTV0NtU1kxT0ZiYmdYdHgwRmpGN3NURzgrdytqOG1uUTZWUAo3V3FTWjA1M2V3RnhrL1hJWGNOd1dBUUQ5bldnM1dKTXdRQURTRGdLR2N0UVFXOERPd09WCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
      verifyHostname: false
      verifyCertificate: true
  extraLabels:
    foo: bar
---
apiVersion: deckhouse.io/v1alpha1
kind: ClusterLogDestination
metadata:
  name: test-es-dest
spec:
  type: Elasticsearch
  elasticsearch:
    index: "logs-%F"
    endpoint: "http://192.168.1.1:9200"
    tls:
      caFile: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3ekNDQWFzQ0ZDalVzcGp5b29wVmdOcjR0TE5SS2hSWERmQXhNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1CNHgKQ3pBSkJnTlZCQVlUQWxKVk1ROHdEUVlEVlFRRERBWlVaWE4wUTBFd0hoY05NakV3TmpJeU1URTBOakEwV2hjTgpORGd4TVRBM01URTBOakEwV2pBZU1Rc3dDUVlEVlFRR0V3SlNWVEVQTUEwR0ExVUVBd3dHVkdWemRFTkJNSUlCCklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEzbG42U3pWSVR1VndlRFRneXR4TDZOTEMKditaeWc5d1dpVllSVnFjZ2hPU0FQMlhSZTJjTWJpYU5vbk9oZW00NDRka0JFY3d4WWhYZVhBWUE0N1dCSHZRRworWkZLOW9KaUJNZGRpSFpmNWpUV1pDK29KKzZMK0h0R2R4MUs3czNZaDM4aUMyWHRqelU5UUJzZmVCZUpIellZCmVXcm1MdDZpTjZRdDQ0Y3l3UHRKVW93ampKaU9YUHYxejluVDdjL3NGLzlTMUVsWENMV1B5dHdKV1NiMGVEUisKYTFGdmdFS1dxTWFySnJFbTFpWVhLU1FZUGFqWE9UU2hHaW9ITVZDK2VzMW55cHN6TG93ZUJ1Vjc5SS9WVnY0YQpnVk5CYTcwaWJEcXM3L3czcTJ3Q2I1ZlpBREU4MzJTcldIdGNtL0luSkNrQUt5czBySTlmODlQWHlHb1lNd0lECkFRQUJNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUM0b3lqL3V0VlFZa242eXU1UTBNbmVPK1YvTlNFSHhqTnIKcldOZnJuT2NTV2I4akFRWjN2ZFpHS1FMVW9raGFTUUNKQndhckxiQWlOVW1udG9ndEhEbEtnZUdxdGdVN3hWeQpJaTFCSlc1VkxyejhMNUdNREdQR2NmUjNpVDdKaDVyelM1UUc5YXlzVHgvMGpWaFN0T1I1cnFqdDlocmZrK0kvClQrT01QTTVrbHpzYXlnZTlkSEx1K3l1VzBzeHhHUk83KzlPeVY3bk9KNEd0TEhicWV0ajBWQUIraWpDMHp1NU0KakxDdm9aZEpQUFViWmVRenFlVW5ZTUwrQ0NERXpCSkdJRk9Xd2w1M2VTblFXbFdVaVJPZWNhd0hobkJzMWlHYgpTQ1BEMTFNMzRRRWZYMHBqQ054RUlzTUtvdFR6V2hFaCsvb0tyQnl2dW16SmpWeWtyU2l5Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
      crtFile: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN0akNDQVo0Q0ZHWDNFQ3I0V3dvVlBhUFpDNGZab042c2JYY09NQTBHQ1NxR1NJYjNEUUVCQ3dVQU1CNHgKQ3pBSkJnTlZCQVlUQWxKVk1ROHdEUVlEVlFRRERBWlVaWE4wUTBFd0hoY05NakV3TmpJeU1URTFOekUyV2hjTgpNelV3TXpBeE1URTFOekUyV2pBUk1ROHdEUVlEVlFRRERBWjJaV04wYjNJd2dnRWlNQTBHQ1NxR1NJYjNEUUVCCkFRVUFBNElCRHdBd2dnRUtBb0lCQVFER0JkSHBvWC9mQytaUkdFQVZpT2tyeE91b0JIazEyYVNLRldVU2hJSFcKZWowNC9zMUtjZFF5RUxlSlk5YUMxTzVuZ1hzdVpDVUNmS1NWdHE1Y3IySTV6cjRaaXNyM0JZK3JlcVBVYkVlYgpLNFBCdEVROUlibno2RTZMVUt3SitIRTFZamliRUFuRkRlamhSUWp6MHFUNWFYR1lNd0RkK1dGMUZ2YzFlUHkvCjhsZEc3YzNvRmczb0ZiV1p6bm9WQmYzOXh3WWZZdEZ2cGN2NWYwbW1SVmZlempRUk9nblhjT1dGb1F4VWcwSjEKV1FFM0xVSUdYMTBzQVpzdUpwMzVSN0tBL1pIRjZHcjhwemZIUmNRaHZPb2VBY0pPdTZZMFBaMnBwSzBhekt6LwpxeHMrZi9hUUJmc0N0c3V2Ty9HbmIvWWFDM1R3QTJmZXhlKzJBWjZGK1NBVEFnTUJBQUV3RFFZSktvWklodmNOCkFRRUxCUUFEZ2dFQkFFeEhkOUtBdkFZYTB2aG1aU0VkR1g3TnZIajhBWDFPV1VBcXZicHJ3YkZ1QkgyZm5LWCsKTmJGVHZXakpDUDdkem10cHphMVQ5RG1vOTJDNC9sWjk0Vy9Vc0pPRjJjSEFRUHlKdk5TdmJPVEg5YTAzajhCaAppbVJ3Zm0rTHNub3RGS3h3VTRhUCtRSEcrRVB2L0FDMDF3UDVhOWVpMEVZWnJIUXh1dTVsOWdURFdjU3Rra1o5Ci8xdzRFWGdNQ2xZVVdnQ1VHUTYvNy9XTkJONTNjWWZ5aU1QcS9VTmVQZUlhUkJDbXJxbklaUCtTWjVwMzFFUXMKZnIyak1rUUo5bTdqNlhWL0RrZFhTSWwrVmdmaVhRSXJDcVN2UXV3RldwdnBicFRPcFJOclhhNGlrMEJLMG1LaQpiYmkwTFVnbzJTcGJuSGlydGlWeVAvMTBCdWhmM3dISUdHUT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
      keyFile: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBeGdYUjZhRi8zd3ZtVVJoQUZZanBLOFRycUFSNU5kbWtpaFZsRW9TQjFubzlPUDdOClNuSFVNaEMzaVdQV2d0VHVaNEY3TG1RbEFueWtsYmF1WEs5aU9jNitHWXJLOXdXUHEzcWoxR3hIbXl1RHdiUkUKUFNHNTgraE9pMUNzQ2ZoeE5XSTRteEFKeFEzbzRVVUk4OUtrK1dseG1ETUEzZmxoZFJiM05Yajh2L0pYUnUzTgo2QllONkJXMW1jNTZGUVg5L2NjR0gyTFJiNlhMK1g5SnBrVlgzczQwRVRvSjEzRGxoYUVNVklOQ2RWa0JOeTFDCkJsOWRMQUdiTGlhZCtVZXlnUDJSeGVocS9LYzN4MFhFSWJ6cUhnSENUcnVtTkQyZHFhU3RHc3lzLzZzYlBuLzIKa0FYN0FyYkxyenZ4cDIvMkdndDA4QU5uM3NYdnRnR2VoZmtnRXdJREFRQUJBb0lCQURVcXd0MXpteDJMMkY3VgpuLzhvTDFLdElJaVFDdXRHY0VNUzAzeFJUM3NDZndXYWhBd0UyL0JGUk1JQ3FFbWdXaEk0VlpaekZPekNBbjZmCitkaXd6akt2SzZNMy9KNnVRNURLOE1uTCtMM1V4Ujl4QXhGV3lOS1FBT2F1MWtJbkRsNUM3T2ZWT29wSjNjajkKL0JWYTdTaDZBeUhXTDlscFo1MUVlVU5HSkxaMEpadWZCMVFiQVdpME5hRVpIdWFPL1FDWU55Qjh5Tk1PQkd5YQpPOUxtZHlDZk85VC9ZTFpXeC9kQ041WldZckhqVEpaREd3T3lCd1k1QjAzUWFmSitxQU5OSkVTTWV6bnlUdkRKCjk5d2hIQ0lxRjRDaHAwM2Y3Sm5QUXJCSDBIbWNDMW9BZjhMWFg5djEvdzY4Smpld1U3VUhoMzlWcTZ0NGNWZXAKdlh4YVdJRUNnWUVBN2dDTFNTVlJQUXFvRlBBcHhEMDVmQmpNUmd2M2tTbWlwWlVNOW5XMkR2WHNUUlFDVFNTcwpVL2JUMG5xZ0FtVTdXZVI3aUFMM2VKMU5ucjd5alc4ZUxaeXNGWUpvMzJNMmxHUGdIdVZoelJYL3ZuQ05CMUNHCmRrWVh5ZDVyK0grdkk1ZWxIcG8rbFVpYWd2NEtiQmtsQkNnRDllNFd6ZFhXN3F4STljc01PRU1DZ1lFQTFQOVIKeGhGNUJoNGVHV1g3RW1DMFRmMlVDa09wOTF1QXpQZDNmNFNQWHlkS2xxMDJCa3BCeFZKZEN2QVc2WlRGZ3FNdQp0Z1BxRi8rSzRNNy9IRStiODhoNytWdkJNVTIwdHFuNWM1Q2J0TUdlSU04MWkvdWxFODlqUlZ2LzI0Y3hZRitDCmlUdFZwUnh1NElNc05rdnAwNHhCMjZ1cGhHMk5HN0NVY2ZBdEkvRUNnWUVBcmpYQnZvbk5QRFFuc2lQVlBxcGUKQUlNYVN3K0phRDBrcTdVOVpzM2t0SEM0UmZjbWRCY3ErTTdNWDkyWWNBaHZlQzR4YWU1Wi9IU1FFMm5MbTFGQgpzcnRpanVBRktiYXloYzNSaUd2NHVhaW5xVnN6TDY1MnJlNUNqV1g4ZkVuaUJkaURhYklYcXlnWXlWZHdnNDJvCk5iR2dySXhaTHRPZTN0ZEhGSHRLOTRjQ2dZQnFXQ09xNGJSc0lvTmlxUEVuSnRNL0VUbGx1b3pVN0lHdFZHejgKWk9IMFh6aTFiRHZKL2k5Q1pySC9zUW12aTlEbFBiWW51R0tib3NIakpsWm0relJoRGhzZnovandOZHpoU3BJNgphZHZqN3J1Vm8vOFhLZ2dza09IK2trVjNoTk5aUzdadjhBajl5K2xyL1BJSkZmUGo1R1pKV0RibDRKQ1FYNlJ1CkVyMW04UUtCZ0VJdE5JSktDOEtNcjJ4VlBjbmo1NExZZ1BvYnhRcktLTlNnRUMrRTNkRFY4TEQyNnZHSmZRY0kKTDBsUE8zVm1vWWRaQnlraUF0NUNYRzUvRks5SkNTV0NtU1kxT0ZiYmdYdHgwRmpGN3NURzgrdytqOG1uUTZWUAo3V3FTWjA1M2V3RnhrL1hJWGNOd1dBUUQ5bldnM1dKTXdRQURTRGdLR2N0UVFXOERPd09WCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
      verifyHostname: false
    auth:
      strategy: Basic
      user: elastic
      password: c2VjcmV0
  extraLabels:
    foo: bar
---
`, 1))
			f.RunHook()
		})

		It("Should create secret", func() {
			Expect(f).To(ExecuteSuccessfully())

			Expect(f.ValuesGet("logShipper.internal.activated").Bool()).To(BeTrue())
			secret := f.KubernetesResource("Secret", "d8-log-shipper", "d8-log-shipper-config")
			Expect(secret).To(Not(BeEmpty()))
			config := secret.Field(`data`).Get("vector\\.json").String()
			d, _ := base64.StdEncoding.DecodeString(config)
			Expect(d).Should(MatchJSON(`
			{
				"sources": {
				  "d8_cluster_test-source_test-es-dest": {
					"type": "kubernetes_logs",
					"extra_label_selector": "app=test",
					"extra_field_selector": "metadata.namespace=tests-whispers",
					"annotation_fields": {
					  "container_image": "image",
					  "container_name": "container",
					  "pod_ip": "pod_ip",
					  "pod_labels": "pod_labels",
					  "pod_name": "pod",
					  "pod_namespace": "namespace",
					  "pod_node_name": "node",
					  "pod_owner": "pod_owner"
					}
				  },
				  "d8_cluster_test-source_test-logstash-dest": {
					"type": "kubernetes_logs",
					"extra_label_selector": "app=test",
					"extra_field_selector": "metadata.namespace=tests-whispers",
					"annotation_fields": {
					  "container_image": "image",
					  "container_name": "container",
					  "pod_ip": "pod_ip",
					  "pod_labels": "pod_labels",
					  "pod_name": "pod",
					  "pod_namespace": "namespace",
					  "pod_node_name": "node",
					  "pod_owner": "pod_owner"
					}
				  },
				  "d8_cluster_test-source_test-loki-dest": {
					"type": "kubernetes_logs",
					"extra_label_selector": "app=test",
					"extra_field_selector": "metadata.namespace=tests-whispers",
					"annotation_fields": {
					  "container_image": "image",
					  "container_name": "container",
					  "pod_ip": "pod_ip",
					  "pod_labels": "pod_labels",
					  "pod_name": "pod",
					  "pod_namespace": "namespace",
					  "pod_node_name": "node",
					  "pod_owner": "pod_owner"
					}
				  }
				},
				"transforms": {
					"d8_tf_test-source_test-es-dest_0": {
					  "group_by": [
						"file",
						"stream"
					  ],
					  "inputs": [
						"d8_cluster_test-source_test-es-dest"
					  ],
					  "merge_strategies": {
						"message": "concat"
					  },
					  "starts_when": " match!(.message, r'^Traceback|^[ ]+|(ERROR|INFO|DEBUG|WARN)') || match!(.message, r'^((([a-zA-Z\\-0-9]+)_([a-zA-Z\\-0-9]+)\\s)|(([a-zA-Z\\-0-9]+)\\s)|(.{0}))(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}) \\[|^(\\{\\s{0,1}\")|^(\\d{2}-\\w{3}-\\d{4}\\s\\d{2}:\\d{2}:\\d{2}\\.{0,1}\\d{2,3})\\s(\\w+)|^([A-Z][0-9]{0,4}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{0,6})') || match!(.message, r'^[^\\s]') ",
					  "type": "reduce"
					},
					"d8_tf_test-source_test-es-dest_1": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_test-source_test-es-dest_0"
					  ],
					  "source": " label1 = .pod_labels.\"controller-revision-hash\" \n if label1 != null { \n   del(.pod_labels.\"controller-revision-hash\") \n } \n label2 = .pod_labels.\"pod-template-hash\" \n if label2 != null { \n   del(.pod_labels.\"pod-template-hash\") \n } \n label3 = .kubernetes \n if label3 != null { \n   del(.kubernetes) \n } \n label4 = .file \n if label4 != null { \n   del(.file) \n } \n",
					  "type": "remap"
					},
					"d8_tf_test-source_test-es-dest_2": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_test-source_test-es-dest_1"
					  ],
					  "source": " .foo=\"bar\" \n",
					  "type": "remap"
					},
					"d8_tf_test-source_test-es-dest_3": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_test-source_test-es-dest_2"
					  ],
					  "source": " structured, err1 = parse_json(.message) \n if err1 == null { \n   .data = structured \n   del(.message) \n } else { \n   .data.message = del(.message)\n } \n",
					  "type": "remap"
					},
					"d8_tf_test-source_test-logstash-dest_0": {
					  "group_by": [
						"file",
						"stream"
					  ],
					  "inputs": [
						"d8_cluster_test-source_test-logstash-dest"
					  ],
					  "merge_strategies": {
						"message": "concat"
					  },
					  "starts_when": " match!(.message, r'^Traceback|^[ ]+|(ERROR|INFO|DEBUG|WARN)') || match!(.message, r'^((([a-zA-Z\\-0-9]+)_([a-zA-Z\\-0-9]+)\\s)|(([a-zA-Z\\-0-9]+)\\s)|(.{0}))(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}) \\[|^(\\{\\s{0,1}\")|^(\\d{2}-\\w{3}-\\d{4}\\s\\d{2}:\\d{2}:\\d{2}\\.{0,1}\\d{2,3})\\s(\\w+)|^([A-Z][0-9]{0,4}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{0,6})') || match!(.message, r'^[^\\s]') ",
					  "type": "reduce"
					},
					"d8_tf_test-source_test-logstash-dest_1": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_test-source_test-logstash-dest_0"
					  ],
					  "source": " label1 = .pod_labels.\"controller-revision-hash\" \n if label1 != null { \n   del(.pod_labels.\"controller-revision-hash\") \n } \n label2 = .pod_labels.\"pod-template-hash\" \n if label2 != null { \n   del(.pod_labels.\"pod-template-hash\") \n } \n label3 = .kubernetes \n if label3 != null { \n   del(.kubernetes) \n } \n label4 = .file \n if label4 != null { \n   del(.file) \n } \n",
					  "type": "remap"
					},
					"d8_tf_test-source_test-logstash-dest_2": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_test-source_test-logstash-dest_1"
					  ],
					  "source": " .foo=\"bar\" \n",
					  "type": "remap"
					},
					"d8_tf_test-source_test-logstash-dest_3": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_test-source_test-logstash-dest_2"
					  ],
					  "source": " structured, err1 = parse_json(.message) \n if err1 == null { \n   .data = structured \n   del(.message) \n } else { \n   .data.message = del(.message)\n } \n",
					  "type": "remap"
					},
					"d8_tf_test-source_test-loki-dest_0": {
					  "group_by": [
						"file",
						"stream"
					  ],
					  "inputs": [
						"d8_cluster_test-source_test-loki-dest"
					  ],
					  "merge_strategies": {
						"message": "concat"
					  },
					  "starts_when": " match!(.message, r'^Traceback|^[ ]+|(ERROR|INFO|DEBUG|WARN)') || match!(.message, r'^((([a-zA-Z\\-0-9]+)_([a-zA-Z\\-0-9]+)\\s)|(([a-zA-Z\\-0-9]+)\\s)|(.{0}))(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}) \\[|^(\\{\\s{0,1}\")|^(\\d{2}-\\w{3}-\\d{4}\\s\\d{2}:\\d{2}:\\d{2}\\.{0,1}\\d{2,3})\\s(\\w+)|^([A-Z][0-9]{0,4}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{0,6})') || match!(.message, r'^[^\\s]') ",
					  "type": "reduce"
					},
					"d8_tf_test-source_test-loki-dest_1": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_test-source_test-loki-dest_0"
					  ],
					  "source": " label1 = .pod_labels.\"controller-revision-hash\" \n if label1 != null { \n   del(.pod_labels.\"controller-revision-hash\") \n } \n label2 = .pod_labels.\"pod-template-hash\" \n if label2 != null { \n   del(.pod_labels.\"pod-template-hash\") \n } \n label3 = .kubernetes \n if label3 != null { \n   del(.kubernetes) \n } \n label4 = .file \n if label4 != null { \n   del(.file) \n } \n",
					  "type": "remap"
					}
				  },
				"sinks": {
				  "d8_cluster_test-es-dest": {
					"type": "elasticsearch",
					"inputs": [
					  "d8_tf_test-source_test-es-dest_3"
					],
					"healthcheck": {
					  "enabled": false
					},
					"buffer": {
					  "max_size": 104857600,
					  "type": "disk"
					},
					"endpoint": "http://192.168.1.1:9200",
					"encoding": {
					  "timestamp_format": "rfc3339"
					},
					"batch": {
					  "max_bytes": 10485760,
					  "timeout_secs": 1
					},
					"auth": {
					  "password": "secret",
					  "strategy": "basic",
					  "user": "elastic"
					},
					"tls": {
					  "ca_file": "-----BEGIN CERTIFICATE-----\nMIICwzCCAasCFCjUspjyoopVgNr4tLNRKhRXDfAxMA0GCSqGSIb3DQEBCwUAMB4x\nCzAJBgNVBAYTAlJVMQ8wDQYDVQQDDAZUZXN0Q0EwHhcNMjEwNjIyMTE0NjA0WhcN\nNDgxMTA3MTE0NjA0WjAeMQswCQYDVQQGEwJSVTEPMA0GA1UEAwwGVGVzdENBMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3ln6SzVITuVweDTgytxL6NLC\nv+Zyg9wWiVYRVqcghOSAP2XRe2cMbiaNonOhem444dkBEcwxYhXeXAYA47WBHvQG\n+ZFK9oJiBMddiHZf5jTWZC+oJ+6L+HtGdx1K7s3Yh38iC2XtjzU9QBsfeBeJHzYY\neWrmLt6iN6Qt44cywPtJUowjjJiOXPv1z9nT7c/sF/9S1ElXCLWPytwJWSb0eDR+\na1FvgEKWqMarJrEm1iYXKSQYPajXOTShGioHMVC+es1nypszLoweBuV79I/VVv4a\ngVNBa70ibDqs7/w3q2wCb5fZADE832SrWHtcm/InJCkAKys0rI9f89PXyGoYMwID\nAQABMA0GCSqGSIb3DQEBCwUAA4IBAQC4oyj/utVQYkn6yu5Q0MneO+V/NSEHxjNr\nrWNfrnOcSWb8jAQZ3vdZGKQLUokhaSQCJBwarLbAiNUmntogtHDlKgeGqtgU7xVy\nIi1BJW5VLrz8L5GMDGPGcfR3iT7Jh5rzS5QG9aysTx/0jVhStOR5rqjt9hrfk+I/\nT+OMPM5klzsayge9dHLu+yuW0sxxGRO7+9OyV7nOJ4GtLHbqetj0VAB+ijC0zu5M\njLCvoZdJPPUbZeQzqeUnYML+CCDEzBJGIFOWwl53eSnQWlWUiROecawHhnBs1iGb\nSCPD11M34QEfX0pjCNxEIsMKotTzWhEh+/oKrByvumzJjVykrSiy\n-----END CERTIFICATE-----\n",
					  "crt_file": "-----BEGIN CERTIFICATE-----\nMIICtjCCAZ4CFGX3ECr4WwoVPaPZC4fZoN6sbXcOMA0GCSqGSIb3DQEBCwUAMB4x\nCzAJBgNVBAYTAlJVMQ8wDQYDVQQDDAZUZXN0Q0EwHhcNMjEwNjIyMTE1NzE2WhcN\nMzUwMzAxMTE1NzE2WjARMQ8wDQYDVQQDDAZ2ZWN0b3IwggEiMA0GCSqGSIb3DQEB\nAQUAA4IBDwAwggEKAoIBAQDGBdHpoX/fC+ZRGEAViOkrxOuoBHk12aSKFWUShIHW\nej04/s1KcdQyELeJY9aC1O5ngXsuZCUCfKSVtq5cr2I5zr4Zisr3BY+reqPUbEeb\nK4PBtEQ9Ibnz6E6LUKwJ+HE1YjibEAnFDejhRQjz0qT5aXGYMwDd+WF1Fvc1ePy/\n8ldG7c3oFg3oFbWZznoVBf39xwYfYtFvpcv5f0mmRVfezjQROgnXcOWFoQxUg0J1\nWQE3LUIGX10sAZsuJp35R7KA/ZHF6Gr8pzfHRcQhvOoeAcJOu6Y0PZ2ppK0azKz/\nqxs+f/aQBfsCtsuvO/Gnb/YaC3TwA2fexe+2AZ6F+SATAgMBAAEwDQYJKoZIhvcN\nAQELBQADggEBAExHd9KAvAYa0vhmZSEdGX7NvHj8AX1OWUAqvbprwbFuBH2fnKX+\nNbFTvWjJCP7dzmtpza1T9Dmo92C4/lZ94W/UsJOF2cHAQPyJvNSvbOTH9a03j8Bh\nimRwfm+LsnotFKxwU4aP+QHG+EPv/AC01wP5a9ei0EYZrHQxuu5l9gTDWcStkkZ9\n/1w4EXgMClYUWgCUGQ6/7/WNBN53cYfyiMPq/UNePeIaRBCmrqnIZP+SZ5p31EQs\nfr2jMkQJ9m7j6XV/DkdXSIl+VgfiXQIrCqSvQuwFWpvpbpTOpRNrXa4ik0BK0mKi\nbbi0LUgo2SpbnHirtiVyP/10Buhf3wHIGGQ=\n-----END CERTIFICATE-----\n",
					  "key_file": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAxgXR6aF/3wvmURhAFYjpK8TrqAR5NdmkihVlEoSB1no9OP7N\nSnHUMhC3iWPWgtTuZ4F7LmQlAnyklbauXK9iOc6+GYrK9wWPq3qj1GxHmyuDwbRE\nPSG58+hOi1CsCfhxNWI4mxAJxQ3o4UUI89Kk+WlxmDMA3flhdRb3NXj8v/JXRu3N\n6BYN6BW1mc56FQX9/ccGH2LRb6XL+X9JpkVX3s40EToJ13DlhaEMVINCdVkBNy1C\nBl9dLAGbLiad+UeygP2Rxehq/Kc3x0XEIbzqHgHCTrumND2dqaStGsys/6sbPn/2\nkAX7ArbLrzvxp2/2Ggt08ANn3sXvtgGehfkgEwIDAQABAoIBADUqwt1zmx2L2F7V\nn/8oL1KtIIiQCutGcEMS03xRT3sCfwWahAwE2/BFRMICqEmgWhI4VZZzFOzCAn6f\n+diwzjKvK6M3/J6uQ5DK8MnL+L3UxR9xAxFWyNKQAOau1kInDl5C7OfVOopJ3cj9\n/BVa7Sh6AyHWL9lpZ51EeUNGJLZ0JZufB1QbAWi0NaEZHuaO/QCYNyB8yNMOBGya\nO9LmdyCfO9T/YLZWx/dCN5ZWYrHjTJZDGwOyBwY5B03QafJ+qANNJESMeznyTvDJ\n99whHCIqF4Chp03f7JnPQrBH0HmcC1oAf8LXX9v1/w68JjewU7UHh39Vq6t4cVep\nvXxaWIECgYEA7gCLSSVRPQqoFPApxD05fBjMRgv3kSmipZUM9nW2DvXsTRQCTSSs\nU/bT0nqgAmU7WeR7iAL3eJ1Nnr7yjW8eLZysFYJo32M2lGPgHuVhzRX/vnCNB1CG\ndkYXyd5r+H+vI5elHpo+lUiagv4KbBklBCgD9e4WzdXW7qxI9csMOEMCgYEA1P9R\nxhF5Bh4eGWX7EmC0Tf2UCkOp91uAzPd3f4SPXydKlq02BkpBxVJdCvAW6ZTFgqMu\ntgPqF/+K4M7/HE+b88h7+VvBMU20tqn5c5CbtMGeIM81i/ulE89jRVv/24cxYF+C\niTtVpRxu4IMsNkvp04xB26uphG2NG7CUcfAtI/ECgYEArjXBvonNPDQnsiPVPqpe\nAIMaSw+JaD0kq7U9Zs3ktHC4RfcmdBcq+M7MX92YcAhveC4xae5Z/HSQE2nLm1FB\nsrtijuAFKbayhc3RiGv4uainqVszL652re5CjWX8fEniBdiDabIXqygYyVdwg42o\nNbGgrIxZLtOe3tdHFHtK94cCgYBqWCOq4bRsIoNiqPEnJtM/ETlluozU7IGtVGz8\nZOH0Xzi1bDvJ/i9CZrH/sQmvi9DlPbYnuGKbosHjJlZm+zRhDhsfz/jwNdzhSpI6\nadvj7ruVo/8XKggskOH+kkV3hNNZS7Zv8Aj9y+lr/PIJFfPj5GZJWDbl4JCQX6Ru\nEr1m8QKBgEItNIJKC8KMr2xVPcnj54LYgPobxQrKKNSgEC+E3dDV8LD26vGJfQcI\nL0lPO3VmoYdZBykiAt5CXG5/FK9JCSWCmSY1OFbbgXtx0FjF7sTG8+w+j8mnQ6VP\n7WqSZ053ewFxk/XIXcNwWAQD9nWg3WJMwQADSDgKGctQQW8DOwOV\n-----END RSA PRIVATE KEY-----\n",
					  "verify_hostname": false
					},
					"compression": "gzip",
					"bulk_action": "index",
					"index": "logs-%F"
				  },
				  "d8_cluster_test-logstash-dest": {
					"type": "socket",
					"inputs": [
					  "d8_tf_test-source_test-logstash-dest_3"
					],
					"healthcheck": {
					  "enabled": false
					},
					"buffer": {
					  "max_size": 104857600,
					  "type": "disk"
					},
					"address": "192.168.199.252:9009",
					"encoding": {
					  "codec": "json",
					  "timestamp_format": "rfc3339"
					},
					"mode": "tcp",
					"tls": {
					  "ca_file": "-----BEGIN CERTIFICATE-----\nMIICwzCCAasCFCjUspjyoopVgNr4tLNRKhRXDfAxMA0GCSqGSIb3DQEBCwUAMB4x\nCzAJBgNVBAYTAlJVMQ8wDQYDVQQDDAZUZXN0Q0EwHhcNMjEwNjIyMTE0NjA0WhcN\nNDgxMTA3MTE0NjA0WjAeMQswCQYDVQQGEwJSVTEPMA0GA1UEAwwGVGVzdENBMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3ln6SzVITuVweDTgytxL6NLC\nv+Zyg9wWiVYRVqcghOSAP2XRe2cMbiaNonOhem444dkBEcwxYhXeXAYA47WBHvQG\n+ZFK9oJiBMddiHZf5jTWZC+oJ+6L+HtGdx1K7s3Yh38iC2XtjzU9QBsfeBeJHzYY\neWrmLt6iN6Qt44cywPtJUowjjJiOXPv1z9nT7c/sF/9S1ElXCLWPytwJWSb0eDR+\na1FvgEKWqMarJrEm1iYXKSQYPajXOTShGioHMVC+es1nypszLoweBuV79I/VVv4a\ngVNBa70ibDqs7/w3q2wCb5fZADE832SrWHtcm/InJCkAKys0rI9f89PXyGoYMwID\nAQABMA0GCSqGSIb3DQEBCwUAA4IBAQC4oyj/utVQYkn6yu5Q0MneO+V/NSEHxjNr\nrWNfrnOcSWb8jAQZ3vdZGKQLUokhaSQCJBwarLbAiNUmntogtHDlKgeGqtgU7xVy\nIi1BJW5VLrz8L5GMDGPGcfR3iT7Jh5rzS5QG9aysTx/0jVhStOR5rqjt9hrfk+I/\nT+OMPM5klzsayge9dHLu+yuW0sxxGRO7+9OyV7nOJ4GtLHbqetj0VAB+ijC0zu5M\njLCvoZdJPPUbZeQzqeUnYML+CCDEzBJGIFOWwl53eSnQWlWUiROecawHhnBs1iGb\nSCPD11M34QEfX0pjCNxEIsMKotTzWhEh+/oKrByvumzJjVykrSiy\n-----END CERTIFICATE-----\n",
					  "crt_file": "-----BEGIN CERTIFICATE-----\nMIICtjCCAZ4CFGX3ECr4WwoVPaPZC4fZoN6sbXcOMA0GCSqGSIb3DQEBCwUAMB4x\nCzAJBgNVBAYTAlJVMQ8wDQYDVQQDDAZUZXN0Q0EwHhcNMjEwNjIyMTE1NzE2WhcN\nMzUwMzAxMTE1NzE2WjARMQ8wDQYDVQQDDAZ2ZWN0b3IwggEiMA0GCSqGSIb3DQEB\nAQUAA4IBDwAwggEKAoIBAQDGBdHpoX/fC+ZRGEAViOkrxOuoBHk12aSKFWUShIHW\nej04/s1KcdQyELeJY9aC1O5ngXsuZCUCfKSVtq5cr2I5zr4Zisr3BY+reqPUbEeb\nK4PBtEQ9Ibnz6E6LUKwJ+HE1YjibEAnFDejhRQjz0qT5aXGYMwDd+WF1Fvc1ePy/\n8ldG7c3oFg3oFbWZznoVBf39xwYfYtFvpcv5f0mmRVfezjQROgnXcOWFoQxUg0J1\nWQE3LUIGX10sAZsuJp35R7KA/ZHF6Gr8pzfHRcQhvOoeAcJOu6Y0PZ2ppK0azKz/\nqxs+f/aQBfsCtsuvO/Gnb/YaC3TwA2fexe+2AZ6F+SATAgMBAAEwDQYJKoZIhvcN\nAQELBQADggEBAExHd9KAvAYa0vhmZSEdGX7NvHj8AX1OWUAqvbprwbFuBH2fnKX+\nNbFTvWjJCP7dzmtpza1T9Dmo92C4/lZ94W/UsJOF2cHAQPyJvNSvbOTH9a03j8Bh\nimRwfm+LsnotFKxwU4aP+QHG+EPv/AC01wP5a9ei0EYZrHQxuu5l9gTDWcStkkZ9\n/1w4EXgMClYUWgCUGQ6/7/WNBN53cYfyiMPq/UNePeIaRBCmrqnIZP+SZ5p31EQs\nfr2jMkQJ9m7j6XV/DkdXSIl+VgfiXQIrCqSvQuwFWpvpbpTOpRNrXa4ik0BK0mKi\nbbi0LUgo2SpbnHirtiVyP/10Buhf3wHIGGQ=\n-----END CERTIFICATE-----\n",
					  "key_file": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAxgXR6aF/3wvmURhAFYjpK8TrqAR5NdmkihVlEoSB1no9OP7N\nSnHUMhC3iWPWgtTuZ4F7LmQlAnyklbauXK9iOc6+GYrK9wWPq3qj1GxHmyuDwbRE\nPSG58+hOi1CsCfhxNWI4mxAJxQ3o4UUI89Kk+WlxmDMA3flhdRb3NXj8v/JXRu3N\n6BYN6BW1mc56FQX9/ccGH2LRb6XL+X9JpkVX3s40EToJ13DlhaEMVINCdVkBNy1C\nBl9dLAGbLiad+UeygP2Rxehq/Kc3x0XEIbzqHgHCTrumND2dqaStGsys/6sbPn/2\nkAX7ArbLrzvxp2/2Ggt08ANn3sXvtgGehfkgEwIDAQABAoIBADUqwt1zmx2L2F7V\nn/8oL1KtIIiQCutGcEMS03xRT3sCfwWahAwE2/BFRMICqEmgWhI4VZZzFOzCAn6f\n+diwzjKvK6M3/J6uQ5DK8MnL+L3UxR9xAxFWyNKQAOau1kInDl5C7OfVOopJ3cj9\n/BVa7Sh6AyHWL9lpZ51EeUNGJLZ0JZufB1QbAWi0NaEZHuaO/QCYNyB8yNMOBGya\nO9LmdyCfO9T/YLZWx/dCN5ZWYrHjTJZDGwOyBwY5B03QafJ+qANNJESMeznyTvDJ\n99whHCIqF4Chp03f7JnPQrBH0HmcC1oAf8LXX9v1/w68JjewU7UHh39Vq6t4cVep\nvXxaWIECgYEA7gCLSSVRPQqoFPApxD05fBjMRgv3kSmipZUM9nW2DvXsTRQCTSSs\nU/bT0nqgAmU7WeR7iAL3eJ1Nnr7yjW8eLZysFYJo32M2lGPgHuVhzRX/vnCNB1CG\ndkYXyd5r+H+vI5elHpo+lUiagv4KbBklBCgD9e4WzdXW7qxI9csMOEMCgYEA1P9R\nxhF5Bh4eGWX7EmC0Tf2UCkOp91uAzPd3f4SPXydKlq02BkpBxVJdCvAW6ZTFgqMu\ntgPqF/+K4M7/HE+b88h7+VvBMU20tqn5c5CbtMGeIM81i/ulE89jRVv/24cxYF+C\niTtVpRxu4IMsNkvp04xB26uphG2NG7CUcfAtI/ECgYEArjXBvonNPDQnsiPVPqpe\nAIMaSw+JaD0kq7U9Zs3ktHC4RfcmdBcq+M7MX92YcAhveC4xae5Z/HSQE2nLm1FB\nsrtijuAFKbayhc3RiGv4uainqVszL652re5CjWX8fEniBdiDabIXqygYyVdwg42o\nNbGgrIxZLtOe3tdHFHtK94cCgYBqWCOq4bRsIoNiqPEnJtM/ETlluozU7IGtVGz8\nZOH0Xzi1bDvJ/i9CZrH/sQmvi9DlPbYnuGKbosHjJlZm+zRhDhsfz/jwNdzhSpI6\nadvj7ruVo/8XKggskOH+kkV3hNNZS7Zv8Aj9y+lr/PIJFfPj5GZJWDbl4JCQX6Ru\nEr1m8QKBgEItNIJKC8KMr2xVPcnj54LYgPobxQrKKNSgEC+E3dDV8LD26vGJfQcI\nL0lPO3VmoYdZBykiAt5CXG5/FK9JCSWCmSY1OFbbgXtx0FjF7sTG8+w+j8mnQ6VP\n7WqSZ053ewFxk/XIXcNwWAQD9nWg3WJMwQADSDgKGctQQW8DOwOV\n-----END RSA PRIVATE KEY-----\n",
					  "verify_hostname": false,
					  "verify_certificate": true,
					  "enabled": true
					}
				  },
				  "d8_cluster_test-loki-dest": {
					"type": "loki",
					"inputs": [
					  "d8_tf_test-source_test-loki-dest_1"
					],
					"healthcheck": {
					  "enabled": false
					},
					"buffer": {
					  "max_size": 104857600,
					  "type": "disk"
					},
					"encoding": {
					  "codec": "text"
					},
					"endpoint": "http://192.168.1.1:9000",
					"labels": {
					  "container": "{{ container }}",
					  "foo": "bar",
					  "image": "{{ image }}",
					  "namespace": "{{ namespace }}",
					  "node": "{{ node_name }}",
					  "pod": "{{ pod }}",
					  "pod_ip": "{{ pod_ip }}",
					  "stream": "{{ stream }}",
					  "pod_labels": "{{ pod_labels }}",
					  "pod_owner": "{{ pod_owner }}"
					},
					"remove_label_fields": true,
					"out_of_order_action": "rewrite_timestamp"
				  }
				}
			  }
`))
		})
		Context("With deleting object", func() {
			BeforeEach(func() {
				f.BindingContexts.Set(f.KubeStateSet(""))
				f.RunHook()
			})
			It("Should delete secret and deactivate module", func() {
				Expect(f).To(ExecuteSuccessfully())
				Expect(f.ValuesGet("logShipper.internal.activated").Bool()).To(BeFalse())
				Expect(f.KubernetesResource("Secret", "d8-log-shipper", "d8-log-shipper-config").Exists()).To(BeFalse())
			})
		})
	})

	Context("Namespaced source", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSetAndWaitForBindingContexts(`
---
apiVersion: deckhouse.io/v1alpha1
kind: PodLoggingConfig
metadata:
  name: whispers-logs
  namespace: tests-whispers
spec:
  labelSelector:
    matchLabels:
      app: test
  clusterDestinationRefs:
    - loki-storage
    - test-es-dest
---
apiVersion: deckhouse.io/v1alpha1
kind: ClusterLogDestination
metadata:
  name: loki-storage
spec:
  type: Loki
  loki:
    endpoint: http://loki.loki:3100
  extraLabels:
    foo: bar
---
apiVersion: deckhouse.io/v1alpha1
kind: ClusterLogDestination
metadata:
  name: test-es-dest
spec:
  type: Elasticsearch
  elasticsearch:
    index: "logs-%F"
    endpoint: "http://192.168.1.1:9200"
  extraLabels:
    foo: bar
---
`, 1))
			f.RunHook()
		})

		It("Should create secret", func() {
			Expect(f).To(ExecuteSuccessfully())

			secret := f.KubernetesResource("Secret", "d8-log-shipper", "d8-log-shipper-config")
			Expect(f.ValuesGet("logShipper.internal.activated").Bool()).To(BeTrue())
			Expect(secret).To(Not(BeEmpty()))
			config := secret.Field(`data`).Get("vector\\.json").String()
			d, _ := base64.StdEncoding.DecodeString(config)
			Expect(d).Should(MatchJSON(`
			{
				"sources": {
				  "d8_namespaced_tests-whispers_whispers-logs_loki-storage": {
					"type": "kubernetes_logs",
					"extra_label_selector": "app=test",
					"extra_field_selector": "metadata.namespace=tests-whispers",
					"annotation_fields": {
					  "container_image": "image",
					  "container_name": "container",
					  "pod_ip": "pod_ip",
					  "pod_labels": "pod_labels",
					  "pod_name": "pod",
					  "pod_namespace": "namespace",
					  "pod_node_name": "node",
					  "pod_owner": "pod_owner"
					}
				  },
				  "d8_namespaced_tests-whispers_whispers-logs_test-es-dest": {
					"type": "kubernetes_logs",
					"extra_label_selector": "app=test",
					"extra_field_selector": "metadata.namespace=tests-whispers",
					"annotation_fields": {
					  "container_image": "image",
					  "container_name": "container",
					  "pod_ip": "pod_ip",
					  "pod_labels": "pod_labels",
					  "pod_name": "pod",
					  "pod_namespace": "namespace",
					  "pod_node_name": "node",
					  "pod_owner": "pod_owner"
					}
				  }
				},
				"transforms": {
					"d8_tf_tests-whispers_whispers-logs_loki-storage_0": {
					  "group_by": [
						"file",
						"stream"
					  ],
					  "inputs": [
						"d8_namespaced_tests-whispers_whispers-logs_loki-storage"
					  ],
					  "merge_strategies": {
						"message": "concat"
					  },
					  "starts_when": " match!(.message, r'^Traceback|^[ ]+|(ERROR|INFO|DEBUG|WARN)') || match!(.message, r'^((([a-zA-Z\\-0-9]+)_([a-zA-Z\\-0-9]+)\\s)|(([a-zA-Z\\-0-9]+)\\s)|(.{0}))(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}) \\[|^(\\{\\s{0,1}\")|^(\\d{2}-\\w{3}-\\d{4}\\s\\d{2}:\\d{2}:\\d{2}\\.{0,1}\\d{2,3})\\s(\\w+)|^([A-Z][0-9]{0,4}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{0,6})') || match!(.message, r'^[^\\s]') ",
					  "type": "reduce"
					},
					"d8_tf_tests-whispers_whispers-logs_loki-storage_1": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_tests-whispers_whispers-logs_loki-storage_0"
					  ],
					  "source": " label1 = .pod_labels.\"controller-revision-hash\" \n if label1 != null { \n   del(.pod_labels.\"controller-revision-hash\") \n } \n label2 = .pod_labels.\"pod-template-hash\" \n if label2 != null { \n   del(.pod_labels.\"pod-template-hash\") \n } \n label3 = .kubernetes \n if label3 != null { \n   del(.kubernetes) \n } \n label4 = .file \n if label4 != null { \n   del(.file) \n } \n",
					  "type": "remap"
					},
					"d8_tf_tests-whispers_whispers-logs_test-es-dest_0": {
					  "group_by": [
						"file",
						"stream"
					  ],
					  "inputs": [
						"d8_namespaced_tests-whispers_whispers-logs_test-es-dest"
					  ],
					  "merge_strategies": {
						"message": "concat"
					  },
					  "starts_when": " match!(.message, r'^Traceback|^[ ]+|(ERROR|INFO|DEBUG|WARN)') || match!(.message, r'^((([a-zA-Z\\-0-9]+)_([a-zA-Z\\-0-9]+)\\s)|(([a-zA-Z\\-0-9]+)\\s)|(.{0}))(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}) \\[|^(\\{\\s{0,1}\")|^(\\d{2}-\\w{3}-\\d{4}\\s\\d{2}:\\d{2}:\\d{2}\\.{0,1}\\d{2,3})\\s(\\w+)|^([A-Z][0-9]{0,4}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{0,6})') || match!(.message, r'^[^\\s]') ",
					  "type": "reduce"
					},
					"d8_tf_tests-whispers_whispers-logs_test-es-dest_1": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_tests-whispers_whispers-logs_test-es-dest_0"
					  ],
					  "source": " label1 = .pod_labels.\"controller-revision-hash\" \n if label1 != null { \n   del(.pod_labels.\"controller-revision-hash\") \n } \n label2 = .pod_labels.\"pod-template-hash\" \n if label2 != null { \n   del(.pod_labels.\"pod-template-hash\") \n } \n label3 = .kubernetes \n if label3 != null { \n   del(.kubernetes) \n } \n label4 = .file \n if label4 != null { \n   del(.file) \n } \n",
					  "type": "remap"
					},
					"d8_tf_tests-whispers_whispers-logs_test-es-dest_2": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_tests-whispers_whispers-logs_test-es-dest_1"
					  ],
					  "source": " .foo=\"bar\" \n",
					  "type": "remap"
					},
					"d8_tf_tests-whispers_whispers-logs_test-es-dest_3": {
					  "drop_on_abort": false,
					  "inputs": [
						"d8_tf_tests-whispers_whispers-logs_test-es-dest_2"
					  ],
					  "source": " structured, err1 = parse_json(.message) \n if err1 == null { \n   .data = structured \n   del(.message) \n } else { \n   .data.message = del(.message)\n } \n",
					  "type": "remap"
					}
				  },
				"sinks": {
				  "d8_cluster_loki-storage": {
					"type": "loki",
					"inputs": [
					  "d8_tf_tests-whispers_whispers-logs_loki-storage_1"
					],
					"healthcheck": {
					  "enabled": false
					},
					"buffer": {
					  "max_size": 104857600,
					  "type": "disk"
					},
					"encoding": {
					  "codec": "text"
					},
					"endpoint": "http://loki.loki:3100",
					"labels": {
					  "container": "{{ container }}",
					  "image": "{{ image }}",
					  "namespace": "{{ namespace }}",
					  "node": "{{ node_name }}",
					  "pod": "{{ pod }}",
					  "pod_ip": "{{ pod_ip }}",
					  "stream": "{{ stream }}",
					  "foo": "bar",
					  "pod_labels": "{{ pod_labels }}",
					  "pod_owner": "{{ pod_owner }}"
					},
					"remove_label_fields": true,
					"out_of_order_action": "rewrite_timestamp"
				  },
				  "d8_cluster_test-es-dest": {
					"type": "elasticsearch",
					"inputs": [
					  "d8_tf_tests-whispers_whispers-logs_test-es-dest_3"
					],
					"healthcheck": {
					  "enabled": false
					},
					"buffer": {
					  "max_size": 104857600,
					  "type": "disk"
					},
					"endpoint": "http://192.168.1.1:9200",
					"encoding": {
					  "timestamp_format": "rfc3339"
					},
					"batch": {
					  "max_bytes": 10485760,
					  "timeout_secs": 1
					},
					"compression": "gzip",
					"bulk_action": "index",
					"index": "logs-%F"
				  }
				}
			  }
`))
		})
		Context("With deleting object", func() {
			BeforeEach(func() {
				f.BindingContexts.Set(f.KubeStateSet(""))
				f.RunHook()
			})
			It("Should delete secret and deactivate module", func() {
				Expect(f).To(ExecuteSuccessfully())
				Expect(f.ValuesGet("logShipper.internal.activated").Bool()).To(BeFalse())
				Expect(f.KubernetesResource("Secret", "d8-log-shipper", "d8-log-shipper-config").Exists()).To(BeFalse())
			})
		})
	})

})
