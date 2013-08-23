# Maintainer: AlphaBernd <whoami@dev-urandom.eu>

pkgname=gbt
pkgver=20130823
pkgrel=2
pkgdesc="IRC bot written in go"
arch=('x86_64' 'i686')
url="http://gbt.dev-urandom.eu/"
license=('custom:pizzaware')
makedepends=('git' 'go')
options=('!strip' '!emptydirs')
_gourl=github.com/krautchan/gbt

build() {
  cd "$srcdir"
  GOPATH="$srcdir" go get -d ${_gourl}
  GOPATH="$srcdir" go install -v -ldflags "-X github.com/krautchan/gbt/config.Version $(git --git-dir=$srcdir/src/${_gourl}/.git describe --always)" ${_gourl}
}

package() {
  source /etc/profile.d/go.sh
  mkdir -p "$pkgdir/$GOPATH"
  cp -Rv --preserve=timestamps ${srcdir}/{src,pkg} "$pkgdir/$GOPATH"
  
  strip "${srcdir}/bin/$pkgname"
  cp -Rv --preserve=timestamps ${srcdir}/bin "$pkgdir/$GOROOT"

  # Package license (if available)
  for f in LICENSE COPYING LICENSE.* COPYING.*; do
    if [ -e "$srcdir/src/$_gourl/$f" ]; then
      install -Dm644 "$srcdir/src/$_gourl/$f" \
        "$pkgdir/usr/share/licenses/$pkgname/$f"
    fi
  done
}

# vim:set ts=2 sw=2 et:
