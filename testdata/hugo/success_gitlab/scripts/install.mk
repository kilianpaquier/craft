# Code generated by craft; DO NOT EDIT.

define install_hugo
current_version=$(hugo version || echo "hugo v0.0.0")
new_version=$(curl -fsSL "https://api.github.com/repos/gohugoio/hugo/releases/latest" | jq -r '.tag_name')
if echo "${current_version}" | grep -Eq "${new_version}"; then
	echo "latest hugo version already installed"
	exit 0
fi

echo "installing latest hugo version"
rm -rf "${HOME}/.local/hugo" && mkdir -p "${HOME}/.local/hugo"
curl -fsSL "https://github.com/gohugoio/hugo/releases/download/${new_version}/hugo_extended_${new_version#v*}_linux-amd64.tar.gz" | (cd "${HOME}/.local/hugo" && tar -xz)
chmod +x "${HOME}/.local/hugo/hugo" && ln -sf "${HOME}/.local/hugo/hugo" "${HOME}/.local/bin/hugo"
endef
.PHONY: install-hugo
install-hugo: ; @$(value install_hugo)
.ONESHELL:

define install_netlify
if which netlify >/dev/null; then
	echo "netlify already installed"
	exit 0
fi

echo "installing netlify"
for bin in "bun" "pnpm" "npm"; do
	if which ${bin}; then
		${bin} i -g netlify-cli
		exit $?
	fi
done
endef
.PHONY: install-netlify
install-netlify: ; @$(value install_netlify)
.ONESHELL: