#/usr/bin/env bash

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

function main() {
    pushd "$(dirname "${script_dir}")" > /dev/null
    bump_go_mod_deps
    popd > /dev/null
}

function bump_go_mod_deps() {
    local old_IFS  dep  dep_name  latest_version  sed_opts

    old_IFS="${IFS}"
    IFS=$'\n'
    for dep in $(go list -m -u -f '{{if and (not .Indirect) .Update }}{{.}}{{end}}' all 2> /dev/null); do
        dep_name=$(cut -d" " -f1 <<< "${dep}")
        latest_version=$(cut -d" " -f3 <<< "${dep}" | sed -e 's/^\[//; s/\]$//')

        sed_opts=("-i")
        if [[ $(uname -s) == Darwin ]]; then
            sed_opts+=("" "-E") # BSD sed on macOS
        else
            sed_opts+=("-r")    # GNU sed on Linux
        fi
        (
            set -x
            sed "${sed_opts[@]}" -e "s,^([[:space:]])*${dep_name} .*$,\1${dep_name} ${latest_version}," "go.mod"
        )
    done
    IFS="${old_IFS}"
}

main "$@"
