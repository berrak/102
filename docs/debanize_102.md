## Debanize 102
I's not much debanize diffrence despite the fact that we have multiple source files, but for completeness lets continue.

### Go version
Actual host is Debian stretch, but with stretch-backport. 
```bash
go version
go version go1.8.1 linux/amd64
```

### Use dh-make-golang
Set environment
```bash
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$GOBIN
```
Use option type *program|library* since dh-make-golang can easily be confused by the files in the repository.
```bash
mkdir -p ~/PKGDEB/102
cd ~/PKGDEB/102
dh-make-golang -type program github.com/berrak/102
```
Directory (102) and files after above run:
```bash
102
102_0.0~git20180303.53e3bf5.orig.tar.xz
itp-102.txt
```
Binary application name is not changed by Debian, thus it remains as *102*.
Ignore the text file, its only for Debian internal usage. Lets look into **102** directory.
It contains three git branches:
```bash
git vbr
* master
  pristine-tar
  upstream
```
Switch to pristine-tar branch and confirm have same base name as the original tarball:
```bash
git sw2 pristine-tar
ls -l
102_0.0~git20180303.53e3bf5.orig.tar.xz.delta
102_0.0~git20180303.53e3bf5.orig.tar.xz.id

git sw2 master
```
Edit **changelog** in the debian directory.
```bash
cd debian
vi changelog
```
Replace *UNRELEASED* with stable and remove TODO with (Closes: #123456) or whatever.
Next view the file **control**, but initially it is acceptable for the first packaging attempt.
Edit **rules** and add DH_GOPKG stanza as below:
```bash
#!/usr/bin/make -f

export DH_GOPKG := github.com/berrak/102

%:
        dh $@ --buildsystem=golang --with=golang    
```
Now we will try to build it first with fakeroot to identify problems with these files or the repository:
```bash
cd ..
fakeroot debian/rules build
dh build --buildsystem=golang --with=golang
   dh_testdir -O--buildsystem=golang
   dh_update_autotools_config -O--buildsystem=golang
   dh_autoreconf -O--buildsystem=golang
   dh_auto_configure -O--buildsystem=golang
   dh_auto_build -O--buildsystem=golang
	go install -v -p 1 github.com/berrak/102
github.com/berrak/102
   dh_auto_test -O--buildsystem=golang
	go test -v -p 1 github.com/berrak/102
=== RUN   TestMain
--- PASS: TestMain (0.00s)
=== RUN   TestOne
--- PASS: TestOne (0.00s)
=== RUN   TestTwo
--- PASS: TestTwo (0.00s)
=== RUN   TestThree
--- PASS: TestThree (0.00s)
PASS
ok  	github.com/berrak/102	0.004s
   create-stamp debian/debhelper-build-stamp
```
Nice, it looks promising. Try to build the binary.
```bash
fakeroot debian/rules binary
dh binary --buildsystem=golang --with=golang
   dh_testroot -O--buildsystem=golang
   dh_prep -O--buildsystem=golang
   dh_auto_install -O--buildsystem=golang
	mkdir -p /home/bekr/PKGDEB/102/102/debian/102/usr/share/gocode/src/github.com/berrak/102
	cp -r -T src/github.com/berrak/102 /home/bekr/PKGDEB/102/102/debian/102/usr/share/gocode/src/github.com/berrak/102
   dh_installdocs -O--buildsystem=golang
   dh_installchangelogs -O--buildsystem=golang
   dh_perl -O--buildsystem=golang
   dh_link -O--buildsystem=golang
   dh_strip_nondeterminism -O--buildsystem=golang
   dh_compress -O--buildsystem=golang
   dh_fixperms -O--buildsystem=golang
   dh_strip -O--buildsystem=golang
   dh_makeshlibs -O--buildsystem=golang
   dh_shlibdeps -O--buildsystem=golang
   dh_installdeb -O--buildsystem=golang
   dh_golang -O--buildsystem=golang
   dh_gencontrol -O--buildsystem=golang
dpkg-gencontrol: warning: Depends field of package 102: unknown substitution variable ${shlibs:Depends}
   dh_md5sums -O--buildsystem=golang
   dh_builddeb -u-Zxz -O--buildsystem=golang
dpkg-deb: building package '102' in '../102_0.0~git20180303.53e3bf5-1_amd64.deb'.
```
We have our debian package built, albeit with a warning about ${shlibs:Depends} from the control file.
We will ignore that for now. Next step is to debanize in a chrooted environment. Set it up with:
```bash
DIST=stretch ARCH=amd64 git-pbuilder create
sudo ln -s /var/cache/pbuilder/base-stretch-amd64.cow base.cow
```
Before removing the build directory **102**, note that no binary is packaged.
```bash
cd debian
tree 102
102
├── DEBIAN
│   ├── control
│   └── md5sums
└── usr
    └── share
        ├── doc
        │   └── 102
        │       ├── changelog.Debian.gz
        │       └── copyright
        └── gocode
            └── src
                └── github.com
                    └── berrak
                        └── 102
                            ├── 102.go
                            ├── 102_test.go
                            ├── one.go
                            ├── onetwothree_test.go
                            ├── three.go
                            └── two.go

10 directories, 10 files
```

We will have to add a directive to include our binary file in the debin directory to take care of that.
Remove the the directory **102** and all new files in **debian** that above runs have created.
Also remove the directory ../*obj-x86_64-linux-gnu*.

```bash
rm -fr 102
rm -fr ../obj-x86_64-linux-gnu
rm 102.substvars
rm debhelper-build-stamp
rm files
```

Before we can build in the chroot we have to update the **rules**.
The binaries should be end up in /opt/ZUL/bin after our enterprise name **ZUL**:

```bash
#!/usr/bin/make -f

TMP  = $(CURDIR)/debian/tmp

GO_SRC := 102.go one.go two.go three.go

export DH_GOPKG := github.com/berrak/102

%:
        dh $@ --buildsystem=golang --with=golang
        
override_dh_auto_build:
        go build $(GO_SRC)

override_dh_auto_install:
        mkdir -p $(TMP)/opt/ZUL/bin
        cp 102 $(TMP)/opt/ZUL/bin        
```
Since we want to have the binary to end up in a new directory,
we have to add a new file *102.install* in the debian directory:
```bash
/opt/ZUL/bin
```
Add all debian files to git:
```bash
cd debian
git add .
git com -m 'Initial packaging'
```
Now we can debanize the 102 application in the chrooted stretch (note: golang-1.7-go (1.7.4-2)) environment:
```bash
gbp buildpackage --git-pbuilder --git-compression=xz
```
A number of new files have been created at parent directory:
```bash
102_0.0~git20180303.53e3bf5-1_amd64.build
102_0.0~git20180303.53e3bf5-1_amd64.buildinfo
102_0.0~git20180303.53e3bf5-1_amd64.changes
102_0.0~git20180303.53e3bf5-1_amd64.deb
102_0.0~git20180303.53e3bf5-1.debian.tar.xz
102_0.0~git20180303.53e3bf5-1.dsc
102_0.0~git20180303.53e3bf5.orig.tar.xz
itp-102.txt
```
Install the 102 application:
```bash
sudo dpkg -i 102_0.0~git20180303.53e3bf5-1_amd64.deb
Selecting previously unselected package 102.
(Reading database ... 158213 files and directories currently installed.)
Preparing to unpack 102_0.0~git20180303.53e3bf5-1_amd64.deb ...
Unpacking 102 (0.0~git20180303.53e3bf5-1) ...
Setting up 102 (0.0~git20180303.53e3bf5-1) ...
```
Run application:
```bash
/opt/ZUL/bin/102
Hello 123
```
To see what is in the package:
```bash
dpkg -L 102
/.
/opt
/opt/ZUL
/opt/ZUL/bin
/opt/ZUL/bin/102
/usr
/usr/share
/usr/share/doc
/usr/share/doc/102
/usr/share/doc/102/changelog.Debian.gz
/usr/share/doc/102/copyright
```

In another chapter in the ZUL enterprise golang journey, usage of all this will be explained.
The short git aliases used here is in user ~/.gitconfig:
```bash
[aliases]
com = commit
sw2 = checkout
vbr = branch -a
```


