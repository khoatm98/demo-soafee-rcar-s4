#!/bin/sh
# --------------------------------------------------------------------------
# |              _    _ _______     .----.      _____         _____        |
# |         /\  | |  | |__   __|  .  ____ .    / ____|  /\   |  __ \       |
# |        /  \ | |  | |  | |    .  / __ \ .  | (___   /  \  | |__) |      |
# |       / /\ \| |  | |  | |   .  / / / / v   \___ \ / /\ \ |  _  /       |
# |      / /__\ \ |__| |  | |   . / /_/ /  .   ____) / /__\ \| | \ \       |
# |     /________\____/   |_|   ^ \____/  .   |_____/________\_|  \_\      |
# |                              . _ _  .                                  |
# --------------------------------------------------------------------------
#
# All Rights Reserved.
# Any use of this source code is subject to a license agreement with the
# AUTOSAR development cooperation.
# More information is available at www.autosar.org.
#
# Disclaimer
#
# This work (specification and/or software implementation) and the material
# contained in it, as released by AUTOSAR, is for the purpose of information
# only. AUTOSAR and the companies that have contributed to it shall not be
# liable for any use of the work.
#
# The material contained in this work is protected by copyright and other
# types of intellectual property rights. The commercial exploitation of the
# material contained in this work requires a license to such intellectual
# property rights.
#
# This work may be utilized or reproduced without any modification, in any
# form or by any means, for informational purposes only. For any other
# purpose, no part of the work may be utilized or reproduced, in any form
# or by any means, without permission in writing from the publisher.
#
# The work has been developed for automotive applications only. It has
# neither been developed, nor tested for non-automotive applications.
#
# The word AUTOSAR and the AUTOSAR logo are registered trademarks.
# --------------------------------------------------------------------------

# AUTOSAR Adaptive Platform Release R22-11 APD Pre-Release
# --------------------------------------------------------
#
# Fetch the used external sources
# ===============================

echo Cloning meta-opendds into meta-opendds
git clone https://github.com/oci-labs/meta-opendds meta-opendds
git -c advice.detachedHead=false -C meta-opendds checkout 18a72b6eda484ce0b09f6af3abfbca1527c1c9b7
echo Cloning meta-openembedded into meta-openembedded
git clone https://github.com/openembedded/meta-openembedded meta-openembedded
git -c advice.detachedHead=false -C meta-openembedded checkout 7203130ed8b58c0df75cb72222ac2bcf546bce44
echo Cloning meta-renesas into meta-rcar
git clone https://github.com/renesas-rcar/meta-renesas meta-rcar
git -c advice.detachedHead=false -C meta-rcar checkout 810e5439d8282e653d89da937627c76f2f34b4af
echo Cloning poky into poky
git clone https://git.yoctoproject.org/git/poky poky
git -c advice.detachedHead=false -C poky checkout 4ddc26f4e4c71b6981898687e2c2e9ce587d15b3

echo Complete

