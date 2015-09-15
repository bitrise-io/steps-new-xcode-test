package main

import "testing"

func testTimeoutWith(t *testing.T, outputToSearchIn string, isShouldFind bool) {
	if isFound, err := isTimeOutError(outputToSearchIn); err != nil {
		t.Logf("Input string to search in was: %s", outputToSearchIn)
		t.Fatalf("Error: Expected (nil), actual (%s)", err)
	} else if isFound != isShouldFind {
		t.Logf("Input string to search in was: %s", outputToSearchIn)
		t.Fatalf("Expected (%v), actual (%v)", isShouldFind, isFound)
	}
}

func TestIsTimeOutError(t *testing.T) {
	testTimeoutWith(t, "", false)

	// Should found
	testTimeoutWith(t, "iPhoneSimulator: Timed out waiting", true)

	// Should found
	testTimeoutWith(t, "iphoneSimulator: timed out waiting", true)

	// Should found
	testTimeoutWith(t, "iphoneSimulator: timed out waiting, test test test", true)

	// Should not found
	testTimeoutWith(t, "iphoneSimulator:", false)

	// full test - not found
	longOutputStrNotFound := `
    export SEPARATE_SYMBOL_EDIT=NO
    export SET_DIR_MODE_OWNER_GROUP=YES
    export SET_FILE_MODE_OWNER_GROUP=NO
    export SHALLOW_BUNDLE=YES
    export SHARED_DERIVED_FILE_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/DerivedSources
    export SHARED_FRAMEWORKS_FOLDER_PATH=BitriseSampleWithYMLTests.xctest/SharedFrameworks
    export SHARED_PRECOMPS_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/PrecompiledHeaders
    export SHARED_SUPPORT_FOLDER_PATH=BitriseSampleWithYMLTests.xctest/SharedSupport
    export SKIP_INSTALL=YES
    export SOURCE_ROOT=/Users/awesome-bitrise-user/develop/bitrise/bitrise-yml-converter-test
    export SRCROOT=/Users/awesome-bitrise-user/develop/bitrise/bitrise-yml-converter-test
    export STRINGS_FILE_OUTPUT_ENCODING=binary
    export STRIP_INSTALLED_PRODUCT=YES
    export STRIP_STYLE=non-global
    export SUPPORTED_DEVICE_FAMILIES="1 2"
    export SUPPORTED_PLATFORMS="iphonesimulator iphoneos"
    export SWIFT_OPTIMIZATION_LEVEL=-Onone
    export SYMROOT=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products
    export SYSTEM_ADMIN_APPS_DIR=/Applications/Utilities
    export SYSTEM_APPS_DIR=/Applications
    export SYSTEM_CORE_SERVICES_DIR=/System/Library/CoreServices
    export SYSTEM_DEMOS_DIR=/Applications/Extras
    export SYSTEM_DEVELOPER_APPS_DIR=/Applications/Xcode.app/Contents/Developer/Applications
    export SYSTEM_DEVELOPER_BIN_DIR=/Applications/Xcode.app/Contents/Developer/usr/bin
    export SYSTEM_DEVELOPER_DEMOS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Utilities/Built Examples"
    export SYSTEM_DEVELOPER_DIR=/Applications/Xcode.app/Contents/Developer
    export SYSTEM_DEVELOPER_DOC_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library"
    export SYSTEM_DEVELOPER_GRAPHICS_TOOLS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Graphics Tools"
    export SYSTEM_DEVELOPER_JAVA_TOOLS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Java Tools"
    export SYSTEM_DEVELOPER_PERFORMANCE_TOOLS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Performance Tools"
    export SYSTEM_DEVELOPER_RELEASENOTES_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library/releasenotes"
    export SYSTEM_DEVELOPER_TOOLS=/Applications/Xcode.app/Contents/Developer/Tools
    export SYSTEM_DEVELOPER_TOOLS_DOC_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library/documentation/DeveloperTools"
    export SYSTEM_DEVELOPER_TOOLS_RELEASENOTES_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library/releasenotes/DeveloperTools"
    export SYSTEM_DEVELOPER_USR_DIR=/Applications/Xcode.app/Contents/Developer/usr
    export SYSTEM_DEVELOPER_UTILITIES_DIR=/Applications/Xcode.app/Contents/Developer/Applications/Utilities
    export SYSTEM_DOCUMENTATION_DIR=/Library/Documentation
    export SYSTEM_KEXT_INSTALL_PATH=/System/Library/Extensions
    export SYSTEM_LIBRARY_DIR=/System/Library
    export TARGETED_DEVICE_FAMILY=1,2
    export TARGETNAME=BitriseSampleWithYMLTests
    export TARGET_BUILD_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator
    export TARGET_NAME=BitriseSampleWithYMLTests
    export TARGET_TEMP_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_FILES_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_FILE_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_ROOT=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates
    export TEST_FRAMEWORK_SEARCH_PATHS=" /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/Library/Frameworks /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/iPhoneSimulator8.4.sdk/Developer/Library/Frameworks"
    export TEST_HOST=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYML.app/BitriseSampleWithYML
    export TOOLCHAINS=com.apple.dt.toolchain.iOS8_4
    export TREAT_MISSING_BASELINES_AS_TEST_FAILURES=NO
    export UID=501
    export UNLOCALIZED_RESOURCES_FOLDER_PATH=BitriseSampleWithYMLTests.xctest
    export UNSTRIPPED_PRODUCT=NO
    export USER=awesome-bitrise-user
    export USER_APPS_DIR=/Users/awesome-bitrise-user/Applications
    export USER_LIBRARY_DIR=/Users/awesome-bitrise-user/Library
    export USE_DYNAMIC_NO_PIC=YES
    export USE_HEADERMAP=YES
    export USE_HEADER_SYMLINKS=NO
    export VALIDATE_PRODUCT=NO
    export VALID_ARCHS="i386 x86_64"
    export VERBOSE_PBXCP=NO
    export VERSIONPLIST_PATH=BitriseSampleWithYMLTests.xctest/version.plist
    export VERSION_INFO_BUILDER=awesome-bitrise-user
    export VERSION_INFO_FILE=BitriseSampleWithYMLTests_vers.c
    export VERSION_INFO_STRING="\"@(#)PROGRAM:BitriseSampleWithYMLTests  PROJECT:BitriseSampleWithYML-\""
    export WRAPPER_EXTENSION=xctest
    export WRAPPER_NAME=BitriseSampleWithYMLTests.xctest
    export WRAPPER_SUFFIX=.xctest
    export XCODE_APP_SUPPORT_DIR=/Applications/Xcode.app/Contents/Developer/Library/Xcode
    export XCODE_PRODUCT_BUILD_VERSION=6E35b
    export XCODE_VERSION_ACTUAL=0640
    export XCODE_VERSION_MAJOR=0600
    export XCODE_VERSION_MINOR=0640
    export XPCSERVICES_FOLDER_PATH=BitriseSampleWithYMLTests.xctest/XPCServices
    export YACC=yacc
    export arch=x86_64
    export variant=normal
    /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/swift-stdlib-tool --verbose --copy
Copying libswiftCore.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftCoreGraphics.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftFoundation.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftSecurity.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftUIKit.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftXCTest.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftDispatch.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftObjectiveC.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftCoreImage.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftDarwin.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks

Touch /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest
    cd /Users/awesome-bitrise-user/develop/bitrise/bitrise-yml-converter-test
    export PATH="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/usr/bin:/Applications/Xcode.app/Contents/Developer/usr/bin:/Users/awesome-bitrise-user/.rbenv/shims:/usr/local/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Users/awesome-bitrise-user/develop/go/bin"
    /usr/bin/touch -c /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest

** TEST SUCCEEDED **`

	testTimeoutWith(t, longOutputStrNotFound, false)

	// full test - FOUND
	longOutputStrDidFind := `
    export SEPARATE_SYMBOL_EDIT=NO
    export SET_DIR_MODE_OWNER_GROUP=YES
    export SET_FILE_MODE_OWNER_GROUP=NO
    export SHALLOW_BUNDLE=YES
    export SHARED_DERIVED_FILE_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/DerivedSources
    export SHARED_FRAMEWORKS_FOLDER_PATH=BitriseSampleWithYMLTests.xctest/SharedFrameworks
    export SHARED_PRECOMPS_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/PrecompiledHeaders
    export SHARED_SUPPORT_FOLDER_PATH=BitriseSampleWithYMLTests.xctest/SharedSupport
    export SKIP_INSTALL=YES
    export SOURCE_ROOT=/Users/awesome-bitrise-user/develop/bitrise/bitrise-yml-converter-test
    export SRCROOT=/Users/awesome-bitrise-user/develop/bitrise/bitrise-yml-converter-test
    export STRINGS_FILE_OUTPUT_ENCODING=binary
    export STRIP_INSTALLED_PRODUCT=YES
    export STRIP_STYLE=non-global
    export SUPPORTED_DEVICE_FAMILIES="1 2"
    export SUPPORTED_PLATFORMS="iphonesimulator iphoneos"
    export SWIFT_OPTIMIZATION_LEVEL=-Onone
    export SYMROOT=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products
    export SYSTEM_ADMIN_APPS_DIR=/Applications/Utilities
    export SYSTEM_APPS_DIR=/Applications
    export SYSTEM_CORE_SERVICES_DIR=/System/Library/CoreServices
    export SYSTEM_DEMOS_DIR=/Applications/Extras
    export SYSTEM_DEVELOPER_APPS_DIR=/Applications/Xcode.app/Contents/Developer/Applications
    export SYSTEM_DEVELOPER_BIN_DIR=/Applications/Xcode.app/Contents/Developer/usr/bin
    export SYSTEM_DEVELOPER_DEMOS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Utilities/Built Examples"
    export SYSTEM_DEVELOPER_DIR=/Applications/Xcode.app/Contents/Developer
    export SYSTEM_DEVELOPER_DOC_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library"
    export SYSTEM_DEVELOPER_GRAPHICS_TOOLS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Graphics Tools"
    export SYSTEM_DEVELOPER_JAVA_TOOLS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Java Tools"
    export SYSTEM_DEVELOPER_PERFORMANCE_TOOLS_DIR="/Applications/Xcode.app/Contents/Developer/Applications/Performance Tools"
    export SYSTEM_DEVELOPER_RELEASENOTES_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library/releasenotes"
    export SYSTEM_DEVELOPER_TOOLS=/Applications/Xcode.app/Contents/Developer/Tools
    export SYSTEM_DEVELOPER_TOOLS_DOC_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library/documentation/DeveloperTools"
    export SYSTEM_DEVELOPER_TOOLS_RELEASENOTES_DIR="/Applications/Xcode.app/Contents/Developer/ADC Reference Library/releasenotes/DeveloperTools"
    export SYSTEM_DEVELOPER_USR_DIR=/Applications/Xcode.app/Contents/Developer/usr
    export SYSTEM_DEVELOPER_UTILITIES_DIR=/Applications/Xcode.app/Contents/Developer/Applications/Utilities
    export SYSTEM_DOCUMENTATION_DIR=/Library/Documentation
    export SYSTEM_KEXT_INSTALL_PATH=/System/Library/Extensions
    export SYSTEM_LIBRARY_DIR=/System/Library
    export TARGETED_DEVICE_FAMILY=1,2
    export TARGETNAME=BitriseSampleWithYMLTests
    export TARGET_BUILD_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator
    export TARGET_NAME=BitriseSampleWithYMLTests
    export TARGET_TEMP_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_FILES_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_FILE_DIR=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates/BitriseSampleWithYML.build/Debug-iphonesimulator/BitriseSampleWithYMLTests.build
    export TEMP_ROOT=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Intermediates
    export TEST_FRAMEWORK_SEARCH_PATHS=" /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/Library/Frameworks /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/iPhoneSimulator8.4.sdk/Developer/Library/Frameworks"
    export TEST_HOST=/Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYML.app/BitriseSampleWithYML
    export TOOLCHAINS=com.apple.dt.toolchain.iOS8_4
    export TREAT_MISSING_BASELINES_AS_TEST_FAILURES=NO
    export UID=501
    export UNLOCALIZED_RESOURCES_FOLDER_PATH=BitriseSampleWithYMLTests.xctest
    export UNSTRIPPED_PRODUCT=NO
    export USER=awesome-bitrise-user
    export USER_APPS_DIR=/Users/awesome-bitrise-user/Applications
    export USER_LIBRARY_DIR=/Users/awesome-bitrise-user/Library
    export USE_DYNAMIC_NO_PIC=YES
    export USE_HEADERMAP=YES
    export USE_HEADER_SYMLINKS=NO
    export VALIDATE_PRODUCT=NO
    export VALID_ARCHS="i386 x86_64"
    export VERBOSE_PBXCP=NO
    export VERSIONPLIST_PATH=BitriseSampleWithYMLTests.xctest/version.plist
    export VERSION_INFO_BUILDER=awesome-bitrise-user
    export VERSION_INFO_FILE=BitriseSampleWithYMLTests_vers.c
    export VERSION_INFO_STRING="\"@(#)PROGRAM:BitriseSampleWithYMLTests  PROJECT:BitriseSampleWithYML-\""
    export WRAPPER_EXTENSION=xctest
    export WRAPPER_NAME=BitriseSampleWithYMLTests.xctest
    export WRAPPER_SUFFIX=.xctest
    export XCODE_APP_SUPPORT_DIR=/Applications/Xcode.app/Contents/Developer/Library/Xcode
    export XCODE_PRODUCT_BUILD_VERSION=6E35b
    export XCODE_VERSION_ACTUAL=0640
    export XCODE_VERSION_MAJOR=0600
    export XCODE_VERSION_MINOR=0640
    export XPCSERVICES_FOLDER_PATH=BitriseSampleWithYMLTests.xctest/XPCServices
    export YACC=yacc
    export arch=x86_64

    xxxxxxxxxxxxx iPhoneSimulator: Timed out waiting 120 seconds for simulator to boot, current state is 1.


    export variant=normal
    /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/swift-stdlib-tool --verbose --copy
Copying libswiftCore.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftCoreGraphics.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftFoundation.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftSecurity.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftUIKit.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftXCTest.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftDispatch.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftObjectiveC.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftCoreImage.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks
Copying libswiftDarwin.dylib from /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/swift/iphonesimulator to /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest/Frameworks

Touch /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest
    cd /Users/awesome-bitrise-user/develop/bitrise/bitrise-yml-converter-test
    export PATH="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/usr/bin:/Applications/Xcode.app/Contents/Developer/usr/bin:/Users/awesome-bitrise-user/.rbenv/shims:/usr/local/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Users/awesome-bitrise-user/develop/go/bin"
    /usr/bin/touch -c /Users/awesome-bitrise-user/Library/Developer/Xcode/DerivedData/BitriseSampleWithYML-brllnldkzzfqyjeofwhnmfsdyicc/Build/Products/Debug-iphonesimulator/BitriseSampleWithYMLTests.xctest`

	testTimeoutWith(t, longOutputStrDidFind, true)
}
