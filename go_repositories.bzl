"""This module contains a Gazelle-generated go_repositories macro that mirrors //go.mod.

To update this file, please un "make update-go-bazel-deps" from the project's root directory.

It is recommended that the Gazelle-generated BUILD files are regenerated every time this file is
updated. To do so, please run "make update-go-bazel-files" from the project's root directory.
"""

load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    """This Gazelle-generated macro defines repositories for the Go modules in //go.mod."""
    go_repository(
        name = "cat_dario_mergo",
        importpath = "dario.cat/mergo",
        sum = "h1:AGCNq9Evsj31mOgNPcLyXc+4PNABt905YmuqPYYpBWk=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_aclements_go_gg",
        importpath = "github.com/aclements/go-gg",
        sum = "h1:E5Dzlk3akC+T2Zj1LBHgfPK1y8YWgLDnNDRmG+tpSKw=",
        version = "v0.0.0-20170118225347-6dbb4e4fefb0",
    )
    go_repository(
        name = "com_github_aclements_go_moremath",
        importpath = "github.com/aclements/go-moremath",
        sum = "h1:xlwdaKcTNVW4PtpQb8aKA4Pjy0CdJHEqvFbAnvR5m2g=",
        version = "v0.0.0-20210112150236-f10218a38794",
    )
    go_repository(
        name = "com_github_ajstarks_svgo",
        importpath = "github.com/ajstarks/svgo",
        sum = "h1:slYM766cy2nI3BwyRiyQj/Ud48djTMtMebDqepE95rw=",
        version = "v0.0.0-20211024235047-1546f124cd8b",
    )
    go_repository(
        name = "com_github_anmitsu_go_shlex",
        importpath = "github.com/anmitsu/go-shlex",
        sum = "h1:9AeTilPcZAjCFIImctFaOjnTIavg87rW78vTPkQqLI8=",
        version = "v0.0.0-20200514113438-38f4b401e2be",
    )
    go_repository(
        name = "com_github_armon_go_socks5",
        importpath = "github.com/armon/go-socks5",
        sum = "h1:0CwZNZbxp69SHPdPJAN/hZIm0C4OItdklCFmMRWYpio=",
        version = "v0.0.0-20160902184237-e75332964ef5",
    )
    go_repository(
        name = "com_github_bahlo_generic_list_go",
        importpath = "github.com/bahlo/generic-list-go",
        sum = "h1:5sz/EEAK+ls5wF+NeqDpk5+iNdMDXrh3z3nPnH1Wvgk=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_buger_jsonparser",
        importpath = "github.com/buger/jsonparser",
        sum = "h1:2PnMjfWD7wBILjqQbt530v576A/cAbQvEW9gGIpYMUs=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:o7IhLm0Msx3BaB+n3Ag7L8EVlByGnpq14C4YWiu/gL8=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_bwesterb_go_ristretto",
        importpath = "github.com/bwesterb/go-ristretto",
        sum = "h1:1w53tCkGhCQ5djbat3+MH0BAQ5Kfgbt56UZQ/JMzngw=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_bytedance_sonic",
        importpath = "github.com/bytedance/sonic",
        sum = "h1:6iJ6NqdoxCDr6mbY8h18oSO+cShGSMRGCEo7F2h0x8s=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_chenzhuoyu_base64x",
        importpath = "github.com/chenzhuoyu/base64x",
        sum = "h1:qSGYFH7+jGhDF8vLC+iwCD4WpbV1EBDSzWkJODFLams=",
        version = "v0.0.0-20221115062448-fe3a3abad311",
    )
    go_repository(
        name = "com_github_chzyer_logex",
        importpath = "github.com/chzyer/logex",
        sum = "h1:XHDu3E6q+gdHgsdTPH6ImJMIp436vR6MPtH8gP05QzM=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_chzyer_readline",
        importpath = "github.com/chzyer/readline",
        sum = "h1:upd/6fQk4src78LMRzh5vItIt361/o4uq553V8B5sGI=",
        version = "v1.5.1",
    )
    go_repository(
        name = "com_github_chzyer_test",
        importpath = "github.com/chzyer/test",
        sum = "h1:p3BQDXSxOhOG0P9z6/hGnII4LGiEPOYBhs8asl/fC04=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_cloudflare_circl",
        importpath = "github.com/cloudflare/circl",
        sum = "h1:fE/Qz0QdIGqeWfnwq0RE0R7MI51s0M2E4Ga9kq5AEMs=",
        version = "v1.3.3",
    )
    go_repository(
        name = "com_github_containerd_console",
        importpath = "github.com/containerd/console",
        sum = "h1:lIr7SlA5PxZyMV30bDW0MGbiOPXwc63yRuCP0ARubLw=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_cpuguy83_go_md2man_v2",
        importpath = "github.com/cpuguy83/go-md2man/v2",
        sum = "h1:p1EgwI/C7NhT0JmVkwCD2ZBK8j4aeHQX2pMHHBfMQ6w=",
        version = "v2.0.2",
    )
    go_repository(
        name = "com_github_cratonica_2goarray",
        importpath = "github.com/cratonica/2goarray",
        sum = "h1:Wg+722vs7a2zQH5lR9QWYsVbplKeffaQFIs5FTdfNNo=",
        version = "v0.0.0-20190331194516-514510793eaa",
    )
    go_repository(
        name = "com_github_creack_pty",
        importpath = "github.com/creack/pty",
        sum = "h1:uDmaGzcdjhF4i/plgjmEsriH11Y0o7RKapEf/LDaM3w=",
        version = "v1.1.9",
    )
    go_repository(
        name = "com_github_cyphar_filepath_securejoin",
        importpath = "github.com/cyphar/filepath-securejoin",
        sum = "h1:Ugdm7cg7i6ZK6x3xDF1oEu1nfkyfH53EtKeQYTC3kyg=",
        version = "v0.2.4",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_elazarl_goproxy",
        importpath = "github.com/elazarl/goproxy",
        sum = "h1:mATvB/9r/3gvcejNsXKSkQ6lcIaNec2nyfOdlTBR2lU=",
        version = "v0.0.0-20230808193330-2592e75ae04a",
    )
    go_repository(
        name = "com_github_emirpasic_gods",
        importpath = "github.com/emirpasic/gods",
        sum = "h1:FXtiHYKDGKCW2KzwZKx0iC0PQmdlorYgdFG9jPXJ1Bc=",
        version = "v1.18.1",
    )
    go_repository(
        name = "com_github_felixge_fgprof",
        importpath = "github.com/felixge/fgprof",
        sum = "h1:VvyZxILNuCiUCSXtPtYmmtGvb65nqXh2QFWc0Wpf2/g=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_github_gabriel_vasile_mimetype",
        importpath = "github.com/gabriel-vasile/mimetype",
        sum = "h1:w5qFW6JKBz9Y393Y4q372O9A7cUSequkh1Q7OhCmWKU=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_getlantern_context",
        importpath = "github.com/getlantern/context",
        sum = "h1:NRUJuo3v3WGC/g5YiyF790gut6oQr5f3FBI88Wv0dx4=",
        version = "v0.0.0-20190109183933-c447772a6520",
    )
    go_repository(
        name = "com_github_getlantern_errors",
        importpath = "github.com/getlantern/errors",
        sum = "h1:6uJ+sZ/e03gkbqZ0kUG6mfKoqDb4XMAzMIwlajq19So=",
        version = "v0.0.0-20190325191628-abdb3e3e36f7",
    )
    go_repository(
        name = "com_github_getlantern_golog",
        importpath = "github.com/getlantern/golog",
        sum = "h1:guBYzEaLz0Vfc/jv0czrr2z7qyzTOGC9hiQ0VC+hKjk=",
        version = "v0.0.0-20190830074920-4ef2e798c2d7",
    )
    go_repository(
        name = "com_github_getlantern_hex",
        importpath = "github.com/getlantern/hex",
        sum = "h1:micT5vkcr9tOVk1FiH8SWKID8ultN44Z+yzd2y/Vyb0=",
        version = "v0.0.0-20190417191902-c6586a6fe0b7",
    )
    go_repository(
        name = "com_github_getlantern_hidden",
        importpath = "github.com/getlantern/hidden",
        sum = "h1:XYzSdCbkzOC0FDNrgJqGRo8PCMFOBFL9py72DRs7bmc=",
        version = "v0.0.0-20190325191715-f02dbb02be55",
    )
    go_repository(
        name = "com_github_getlantern_ops",
        importpath = "github.com/getlantern/ops",
        sum = "h1:wrYrQttPS8FHIRSlsrcuKazukx/xqO/PpLZzZXsF+EA=",
        version = "v0.0.0-20190325191751-d70cb0d6f85f",
    )
    go_repository(
        name = "com_github_getlantern_systray",
        importpath = "github.com/getlantern/systray",
        sum = "h1:dCEHtfmvkJG7HZ8lS/sLklTH4RKUcIsKrAD9sThoEBE=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_gin_contrib_cors",
        importpath = "github.com/gin-contrib/cors",
        sum = "h1:oJ6gwtUl3lqV0WEIwM/LxPF1QZ5qe2lGWdY2+bz7y0g=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_gin_contrib_sse",
        importpath = "github.com/gin-contrib/sse",
        sum = "h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_gin_gonic_gin",
        importpath = "github.com/gin-gonic/gin",
        sum = "h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_gliderlabs_ssh",
        importpath = "github.com/gliderlabs/ssh",
        sum = "h1:OcaySEmAQJgyYcArR+gGGTHCyE7nvhEMTlYY+Dp8CpY=",
        version = "v0.3.5",
    )
    go_repository(
        name = "com_github_go_fonts_liberation",
        importpath = "github.com/go-fonts/liberation",
        sum = "h1:jAkAWJP4S+OsrPLZM4/eC9iW7CtHy+HBXrEwZXWo5VM=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_go_git_gcfg",
        importpath = "github.com/go-git/gcfg",
        sum = "h1:+zs/tPmkDkHx3U66DAb0lQFJrpS6731Oaa12ikc+DiI=",
        version = "v1.5.1-0.20230307220236-3a3c6141e376",
    )
    go_repository(
        name = "com_github_go_git_go_billy_v5",
        importpath = "github.com/go-git/go-billy/v5",
        sum = "h1:yEY4yhzCDuMGSv83oGxiBotRzhwhNr8VZyphhiu+mTU=",
        version = "v5.5.0",
    )
    go_repository(
        name = "com_github_go_git_go_git_fixtures_v4",
        importpath = "github.com/go-git/go-git-fixtures/v4",
        sum = "h1:eMje31YglSBqCdIqdhKBW8lokaMrL3uTkpGYlE2OOT4=",
        version = "v4.3.2-0.20231010084843-55a94097c399",
    )
    go_repository(
        name = "com_github_go_git_go_git_v5",
        importpath = "github.com/go-git/go-git/v5",
        sum = "h1:XIZc1p+8YzypNr34itUfSvYJcv+eYdTnTvOZ2vD3cA4=",
        version = "v5.11.0",
    )
    go_repository(
        name = "com_github_go_latex_latex",
        importpath = "github.com/go-latex/latex",
        sum = "h1:6zl3BbBhdnMkpSj2YY30qV3gDcVBGtFgVsV3+/i+mKQ=",
        version = "v0.0.0-20210823091927-c0d11ff05a81",
    )
    go_repository(
        name = "com_github_go_pdf_fpdf",
        importpath = "github.com/go-pdf/fpdf",
        sum = "h1:MlgtGIfsdMEEQJr2le6b/HNr1ZlQwxyWr77r2aj2U/8=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_go_playground_assert_v2",
        importpath = "github.com/go-playground/assert/v2",
        sum = "h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=",
        version = "v2.2.0",
    )
    go_repository(
        name = "com_github_go_playground_locales",
        importpath = "github.com/go-playground/locales",
        sum = "h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_github_go_playground_universal_translator",
        importpath = "github.com/go-playground/universal-translator",
        sum = "h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=",
        version = "v0.18.1",
    )
    go_repository(
        name = "com_github_go_playground_validator_v10",
        importpath = "github.com/go-playground/validator/v10",
        sum = "h1:vgvQWe3XCz3gIeFDm/HnTIbj6UGmg/+t63MyGU2n5js=",
        version = "v10.14.0",
    )
    go_repository(
        name = "com_github_go_sql_driver_mysql",
        importpath = "github.com/go-sql-driver/mysql",
        sum = "h1:g24URVg0OFbNUTx9qqY1IRZ9D9z3iPyi5zKhQZpNwpA=",
        version = "v1.4.1",
    )
    go_repository(
        name = "com_github_go_stack_stack",
        importpath = "github.com/go-stack/stack",
        sum = "h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_goccy_go_json",
        importpath = "github.com/goccy/go-json",
        sum = "h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=",
        version = "v0.10.2",
    )
    go_repository(
        name = "com_github_golang_freetype",
        importpath = "github.com/golang/freetype",
        sum = "h1:DACJavvAHhabrF08vX0COfcOBJRhZ8lUbR+ZWIs0Y5g=",
        version = "v0.0.0-20170609003504-e2365dfdc4a0",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:oI5xCqsCo564l8iNU+DwB5epxmsaqB+rhGL0m5jtYqE=",
        version = "v0.0.0-20210331224755-41bb18bfe9da",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_github_gonum_blas",
        importpath = "github.com/gonum/blas",
        sum = "h1:Q0Jsdxl5jbxouNs1TQYt0gxesYMU4VXRbsTlgDloZ50=",
        version = "v0.0.0-20181208220705-f22b278b28ac",
    )
    go_repository(
        name = "com_github_gonum_floats",
        importpath = "github.com/gonum/floats",
        sum = "h1:EvokxLQsaaQjcWVWSV38221VAK7qc2zhaO17bKys/18=",
        version = "v0.0.0-20181209220543-c233463c7e82",
    )
    go_repository(
        name = "com_github_gonum_internal",
        importpath = "github.com/gonum/internal",
        sum = "h1:8jtTdc+Nfj9AR+0soOeia9UZSvYBvETVHZrugUowJ7M=",
        version = "v0.0.0-20181124074243-f884aa714029",
    )
    go_repository(
        name = "com_github_gonum_lapack",
        importpath = "github.com/gonum/lapack",
        sum = "h1:7qnwS9+oeSiOIsiUMajT+0R7HR6hw5NegnKPmn/94oI=",
        version = "v0.0.0-20181123203213-e4cdc5a0bff9",
    )
    go_repository(
        name = "com_github_gonum_matrix",
        importpath = "github.com/gonum/matrix",
        sum = "h1:V2IgdyerlBa/MxaEFRbV5juy/C3MGdj4ePi+g6ePIp4=",
        version = "v0.0.0-20181209220409-c518dec07be9",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_google_gofuzz",
        importpath = "github.com/google/gofuzz",
        sum = "h1:A8PeW59pxE9IoFRqBp37U+mSNaQoZ46F1f0f863XSXw=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_google_pprof",
        importpath = "github.com/google/pprof",
        sum = "h1:1FjCyPC+syAzJ5/2S8fqdZK1R22vvA0J7JZKcuOIQ7Y=",
        version = "v0.0.0-20211214055906-6f57359322fd",
    )
    go_repository(
        name = "com_github_google_s2a_go",
        importpath = "github.com/google/s2a-go",
        sum = "h1:1kZ/sQM3srePvKs3tXAvQzo66XfcReoqFpIpIccE7Oc=",
        version = "v0.1.4",
    )
    go_repository(
        name = "com_github_google_safehtml",
        importpath = "github.com/google/safehtml",
        sum = "h1:ZOt2VXg4x24bW0m2jtzAOkhoXV0iM8vNKc0paByCZqM=",
        version = "v0.0.2",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_googleapis_enterprise_certificate_proxy",
        importpath = "github.com/googleapis/enterprise-certificate-proxy",
        sum = "h1:yk9/cqRKtT9wXZSsRH9aurXEpJX+U6FLtpYTdC3R06k=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:9V9PWXEsWnPpQhu/PeQIkS4eGzMlTLGgt80cUUI8Ki4=",
        version = "v2.11.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_cloudsql_proxy",
        importpath = "github.com/GoogleCloudPlatform/cloudsql-proxy",
        sum = "h1:8F6fjL5iQP6sArGtPuXh0l6hggdcIpAm4ChjVJE4oTs=",
        version = "v0.0.0-20190129172621-c8b1d7a94ddf",
    )
    go_repository(
        name = "com_github_inconshreveable_mousetrap",
        importpath = "github.com/inconshreveable/mousetrap",
        sum = "h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_invopop_jsonschema",
        importpath = "github.com/invopop/jsonschema",
        sum = "h1:6ovsNSuvn9wEQVOyc72aycBMVQFKz7cPdMJn10CvzRI=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_github_jbenet_go_context",
        importpath = "github.com/jbenet/go-context",
        sum = "h1:BQSFePA1RWJOlocH6Fxy8MmwDt+yVQYULKfN0RoTN8A=",
        version = "v0.0.0-20150711004518-d14ea06fba99",
    )
    go_repository(
        name = "com_github_jedib0t_go_pretty_v6",
        importpath = "github.com/jedib0t/go-pretty/v6",
        sum = "h1:ACteMBRrrmm1gMsXe9PSTOClQ63IXDUt03H5U+UV8OU=",
        version = "v6.5.9",
    )
    go_repository(
        name = "com_github_jmorganca_ollama",
        importpath = "github.com/jmorganca/ollama",
        sum = "h1:On87q9B54yYrTy51L4IqOgH5UHRx/yKOQAZbyS8A45E=",
        version = "v0.1.27",
    )
    go_repository(
        name = "com_github_josharian_intern",
        importpath = "github.com/josharian/intern",
        sum = "h1:vlS4z54oSdjm0bgjRigI+G1HpF+tI+9rE5LLzOg8HmY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_json_iterator_go",
        importpath = "github.com/json-iterator/go",
        sum = "h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=",
        version = "v1.1.12",
    )
    go_repository(
        name = "com_github_kevinburke_ssh_config",
        importpath = "github.com/kevinburke/ssh_config",
        sum = "h1:x584FjTGwHzMwvHx18PXxbBVzfnxogHaAReU4gf13a4=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_klauspost_cpuid_v2",
        importpath = "github.com/klauspost/cpuid/v2",
        sum = "h1:acbojRNwl3o09bUq+yDCtZFc1aiwaAAxtcn8YkZXnvk=",
        version = "v2.2.4",
    )
    go_repository(
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_kr_pty",
        importpath = "github.com/kr/pty",
        sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_leodido_go_urn",
        importpath = "github.com/leodido/go-urn",
        sum = "h1:XlAE/cm/ms7TE/VMVoduSpNBoyc2dOxHs5MZSwAN63Q=",
        version = "v1.2.4",
    )
    go_repository(
        name = "com_github_mailru_easyjson",
        importpath = "github.com/mailru/easyjson",
        sum = "h1:UGYAvKxe3sBsEDzO8ZeWOSlIQfWFlxbzLZe7hwFURr0=",
        version = "v0.7.7",
    )
    go_repository(
        name = "com_github_mattn_go_isatty",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:JITubQf0MOLdlGRuRq+jtsDlekdYPia9ZFsB8h/APPA=",
        version = "v0.0.19",
    )
    go_repository(
        name = "com_github_mattn_go_runewidth",
        importpath = "github.com/mattn/go-runewidth",
        sum = "h1:UNAjwbU9l54TA3KzvqLGxwWjHmMgBUVhBiTjelZgg3U=",
        version = "v0.0.15",
    )
    go_repository(
        name = "com_github_mattn_go_sqlite3",
        importpath = "github.com/mattn/go-sqlite3",
        sum = "h1:qZgc/Rwetq+MtyE18WhzjokPD93dNqLGNT3QJuLvBGw=",
        version = "v1.14.14",
    )
    go_repository(
        name = "com_github_microsoft_go_winio",
        importpath = "github.com/Microsoft/go-winio",
        sum = "h1:9/kr64B9VUZrLm5YYwbGtUJnMgqWVOdUAXu6Migciow=",
        version = "v0.6.1",
    )
    go_repository(
        name = "com_github_mmcloughlin_avo",
        importpath = "github.com/mmcloughlin/avo",
        sum = "h1:nAco9/aI9Lg2kiuROBY6BhCI/z0t5jEvJfjWbL8qXLU=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_modern_go_concurrent",
        importpath = "github.com/modern-go/concurrent",
        sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
        version = "v0.0.0-20180306012644-bacd9c7ef1dd",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        importpath = "github.com/modern-go/reflect2",
        sum = "h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_olekukonko_tablewriter",
        importpath = "github.com/olekukonko/tablewriter",
        sum = "h1:P2Ga83D34wi1o9J6Wh1mRuqd4mF/x/lgBS7N7AbDhec=",
        version = "v0.0.5",
    )
    go_repository(
        name = "com_github_onsi_gomega",
        importpath = "github.com/onsi/gomega",
        sum = "h1:naR28SdDFlqrG6kScpT8VWpu1xWY5nJRCF3XaYyBjhI=",
        version = "v1.27.10",
    )
    go_repository(
        name = "com_github_oxtoacart_bpool",
        importpath = "github.com/oxtoacart/bpool",
        sum = "h1:rp5dCmg/yLR3mgFuSOe4oEnDDmGLROTvMragMUXpTQw=",
        version = "v0.0.0-20190530202638-03653db5a59c",
    )
    go_repository(
        name = "com_github_pbnjay_memory",
        importpath = "github.com/pbnjay/memory",
        sum = "h1:onHthvaw9LFnH4t2DcNVpwGmV9E1BkGknEliJkfwQj0=",
        version = "v0.0.0-20210728143218-7b4eea64cf58",
    )
    go_repository(
        name = "com_github_pborman_uuid",
        importpath = "github.com/pborman/uuid",
        sum = "h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:0ctb6s9mE31h0/lhu+J6OPmVeDxJn+kYnJc2jZR9tGQ=",
        version = "v2.0.8",
    )
    go_repository(
        name = "com_github_pjbgf_sha1cd",
        importpath = "github.com/pjbgf/sha1cd",
        sum = "h1:4D5XXmUUBUl/xQ6IjCkEAbqXskkq/4O7LmGn0AqMDs4=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_pkg_diff",
        importpath = "github.com/pkg/diff",
        sum = "h1:aoZm08cpOy4WuID//EZDgcC4zIxODThtZNPirFr42+A=",
        version = "v0.0.0-20210226163009-20ebb0f2a09e",
    )
    go_repository(
        name = "com_github_pkg_errors",
        importpath = "github.com/pkg/errors",
        sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_pkg_profile",
        importpath = "github.com/pkg/profile",
        sum = "h1:hnbDkaNWPCLMO9wGLdBFTIZvzDrDfBM2072E1S9gJkA=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_protonmail_go_crypto",
        importpath = "github.com/ProtonMail/go-crypto",
        sum = "h1:kkhsdkhsCvIsutKu5zLMgWtgh9YxGCNAw8Ad8hjwfYg=",
        version = "v0.0.0-20230828082145-3c4c8a2d2371",
    )
    go_repository(
        name = "com_github_rivo_uniseg",
        importpath = "github.com/rivo/uniseg",
        sum = "h1:S1pD9weZBuJdFmowNwbpi7BJ8TNftyUImj/0WQi72jY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:cWPaGQEPrBb5/AsnsZesgZZ9yb1OQ+GOISoDNXVBh4M=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_github_russross_blackfriday_v2",
        importpath = "github.com/russross/blackfriday/v2",
        sum = "h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_sashabaranov_go_openai",
        importpath = "github.com/sashabaranov/go-openai",
        sum = "h1:+dkuCADSnwXV02YVJkdphY8XD9AyHLUWwk6V7LB6EL8=",
        version = "v1.19.2",
    )
    go_repository(
        name = "com_github_sergi_go_diff",
        importpath = "github.com/sergi/go-diff",
        sum = "h1:we8PVUC3FE2uYfodKH/nBHMSetSfHDR6scGdBi+erh0=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_sirupsen_logrus",
        importpath = "github.com/sirupsen/logrus",
        sum = "h1:trlNQbNUG3OdDrDil03MCb1H2o9nJ1x4/5LYw7byDE0=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_skeema_knownhosts",
        importpath = "github.com/skeema/knownhosts",
        sum = "h1:SHWdIUa82uGZz+F+47k8SY4QhhI291cXCpopT1lK2AQ=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        importpath = "github.com/spf13/cobra",
        sum = "h1:hyqWnYt1ZQShIddO5kBpj3vu05/++x6tJ6dg8EC572I=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:CcVxjf3Q8PM0mHUKJCdn+eZZtm5yQwehR5yeSVQQcUk=",
        version = "v1.8.4",
    )
    go_repository(
        name = "com_github_twitchyliquid64_golang_asm",
        importpath = "github.com/twitchyliquid64/golang-asm",
        sum = "h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=",
        version = "v0.15.1",
    )
    go_repository(
        name = "com_github_ugorji_go",
        importpath = "github.com/ugorji/go",
        sum = "h1:qYhyWUUd6WbiM+C6JZAUkIJt/1WrjzNHY9+KCIjVqTo=",
        version = "v1.2.7",
    )
    go_repository(
        name = "com_github_ugorji_go_codec",
        importpath = "github.com/ugorji/go/codec",
        sum = "h1:BMaWp1Bb6fHwEtbplGBGJ498wD+LKlNSl25MjdZY4dU=",
        version = "v1.2.11",
    )
    go_repository(
        name = "com_github_urfave_cli_v2",
        importpath = "github.com/urfave/cli/v2",
        sum = "h1:8xSQ6szndafKVRmfyeUMxkNUJQMjL1F2zmsZ+qHpfho=",
        version = "v2.27.1",
    )
    go_repository(
        name = "com_github_wk8_go_ordered_map_v2",
        importpath = "github.com/wk8/go-ordered-map/v2",
        sum = "h1:5h/BUHu93oj4gIdvHHHGsScSTMijfx5PeYkE/fJgbpc=",
        version = "v2.1.8",
    )
    go_repository(
        name = "com_github_xanzy_ssh_agent",
        importpath = "github.com/xanzy/ssh-agent",
        sum = "h1:+/15pJfg/RsTxqYcX6fHqOXZwwMP+2VyYWJeWM2qQFM=",
        version = "v0.3.3",
    )
    go_repository(
        name = "com_github_xrash_smetrics",
        importpath = "github.com/xrash/smetrics",
        sum = "h1:bAn7/zixMGCfxrRTfdpNzjtPYqr8smhKouy9mxVdGPU=",
        version = "v0.0.0-20201216005158-039620a65673",
    )
    go_repository(
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:fVcFKWvrslecOb/tg+Cc05dkeYx540o0FuFt3nUVDoE=",
        version = "v1.4.13",
    )
    go_repository(
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:sdFPBr6xG9/wkBbfhmUz/JmZC7X6LavQgcrVINrKiVA=",
        version = "v0.110.2",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:6aKEtlUiwEpJzM001l0yFkpXmUVXaN8W+fbkb2AZNbg=",
        version = "v1.20.1",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:mg4jlk7mCAj6xXp9UJ4fjI9VUI5rubuGBW5aJ7UnBMY=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:+CmB+K0J/33d0zSQ9SlFWUeCCEn5XJA0ZMZ3pHE9u8k=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:6weCgzRvMg7lzuUurI4697AqIRPU1SvzHhynwpW31jI=",
        version = "v1.29.0",
    )
    go_repository(
        name = "ht_sr_git_sbinet_gg",
        importpath = "git.sr.ht/~sbinet/gg",
        sum = "h1:LNhjNn8DerC8f9DHLz6lS0YYul/b602DUxDgGkd/Aik=",
        version = "v0.3.1",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_errgo_v2",
        importpath = "gopkg.in/errgo.v2",
        sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
        version = "v2.1.0",
    )
    go_repository(
        name = "in_gopkg_warnings_v0",
        importpath = "gopkg.in/warnings.v0",
        sum = "h1:wFXVbFY8DY5/xOe1ECiWdKCzZlxgshcYVNkBHstARME=",
        version = "v0.1.2",
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:y73uSU6J157QMP2kn2r30vwW1A2W2WFwSCGnAVxeaD0=",
        version = "v0.24.0",
    )
    go_repository(
        name = "io_rsc_pdf",
        importpath = "rsc.io/pdf",
        sum = "h1:k1MczvYDUvJBe93bYd7wrZLLUEcLZAuF824/I4e5Xr4=",
        version = "v0.1.1",
    )
    go_repository(
        name = "org_golang_google_api",
        importpath = "google.golang.org/api",
        sum = "h1:q4GJq+cAdMAC7XP7njvQ4tvohGLiSlytuL4BQxbIZ+o=",
        version = "v0.126.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
        version = "v1.6.7",
    )
    go_repository(
        name = "org_golang_google_genproto",
        importpath = "google.golang.org/genproto",
        sum = "h1:8DyZCyvI8mE1IdLy/60bS+52xfymkE72wv1asokgtao=",
        version = "v0.0.0-20230530153820-e85fd2cbaebc",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:kVKPf/IiYSBWEWtkIn6wZXwWGCnLKcC8oWfZvXjsGnM=",
        version = "v0.0.0-20230530153820-e85fd2cbaebc",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:XSJ8Vk1SWuNr8S18z1NZSziL0CPIXLCCMDOEFtHBOFc=",
        version = "v0.0.0-20230530153820-e85fd2cbaebc",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:3Oj82/tFSCeUrRTg/5E/7d/W5A1tj6Ky1ABAuZuv5ag=",
        version = "v1.55.0",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:g0LDEJHgrBl9N9r17Ru3sqWhkIx2NB67okBHPwC7hs8=",
        version = "v1.31.0",
    )
    go_repository(
        name = "org_golang_x_arch",
        importpath = "golang.org/x/arch",
        sum = "h1:02VY4/ZcO/gBOH6PUaoiptASxtXU10jazRCP865E97k=",
        version = "v0.3.0",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:ENy+Az/9Y1vSrlrvBSyna3PITt4tiZLf7sgCjZBX7Wo=",
        version = "v0.19.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:m64FZMko/V45gv0bNmrNYoDEq8U5YUhetc9cBWKS1TQ=",
        version = "v0.0.0-20230817173708-d852ddb80c63",
    )
    go_repository(
        name = "org_golang_x_image",
        importpath = "golang.org/x/image",
        sum = "h1:kOELfmgrmJlw4Cdb7g/QGuB3CvDrXbqEIww/pNtNBm8=",
        version = "v0.15.0",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:rmsUpXtvNzj340zd98LZ4KntptpfRHwpFOHG188oHXc=",
        version = "v0.12.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:AQyQV4dYCvJ7vGmJyKki9+PBdyvhkSd8EIx/qb0AYv4=",
        version = "v0.21.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:6m3ZPmLEFdVxKKWnKq4VqZ60gutO35zm+zrAHVmHyDQ=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_golang_x_perf",
        importpath = "golang.org/x/perf",
        sum = "h1:/xigTF9n9L6Plv6RlldYrF6QUT4bDsLrMS9LjEioIB0=",
        version = "v0.0.0-20240208143119-b26761745961",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:ftCYgMx6zT/asHUrPw8BLLscYtGznsLAnjq5RH9P66E=",
        version = "v0.3.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:25cE3gD+tdBA7lp7QfhuV+rJiE9YXTcS3VG1SqssI/Y=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:mkTF7LCd6WGJNL3K1Ad7kwxNfYAW6a8a8QqtMblp/4U=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=",
        version = "v0.14.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:Iey4qkscZuv0VvIt8E0neZjtPVQFSc870HQ448QgEmQ=",
        version = "v0.13.0",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:H2TDz8ibqkAF6YGhCdN3jS9O0/s90v0rJh3X/OLHEUk=",
        version = "v0.0.0-20220907171357-04be3eba64a2",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:f1IJhK4Km5tBJmaiJXtk/PkL4cdVX6J+tGiM187uT5E=",
        version = "v0.11.0",
    )
    go_repository(
        name = "org_gonum_v1_plot",
        importpath = "gonum.org/v1/plot",
        sum = "h1:dnifSs43YJuNMDzB7v8wV64O4ABBHReuAVAoBxqBqS4=",
        version = "v0.10.1",
    )
