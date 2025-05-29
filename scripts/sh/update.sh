#!/bin/sh

log_info() {
  fg="\033[0;34m"
  reset="\033[0m"
  echo "${fg}$1${reset}"
}

cmd=""
for c in craft.dev craft; do
  [ "$cmd" = "" ] || command -v $c > /dev/null 2>&1 || continue
  cmd=$c
done
if [ "$cmd" = "" ]; then
  echo "No craft generator found, exiting"
  exit 2
fi
log_info "Found craft generator named '$cmd'"

dirs=$(find "$HOME/workspaces" -name testdata -prune -o -name .craft -exec dirname {} +;)
for dir in $dirs; do
  log_info "Updating layout of $dir"
  $cmd -d "$dir"
done
unset dir
