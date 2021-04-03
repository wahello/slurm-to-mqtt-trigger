package main

const name = "slurm-to-mqtt-trigger"
const version = "1.0.0-20210403"

const defaultConfigurationFile = "/etc/slurm-to-mqtt-trigger/config.ini"

const versionText = `%s version %s
Copyright (C) 2021 by Andreas Maus <maus@ypbind.de>
This program comes with ABSOLUTELY NO WARRANTY.

pkidb is distributed under the Terms of the GNU General
Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)

Build with go version: %s
`

const helpText = `Usage: %s [--help] [--version]
`
