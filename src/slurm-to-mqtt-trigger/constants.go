package main

const name = "slurm-to-mqtt-trigger"
const version = "1.0.1"

const defaultConfigurationFile = "/etc/slurm-to-mqtt-trigger/config.ini"

const versionText = `%s version %s
Copyright (C) 2021 by Andreas Maus <maus@ypbind.de>
This program comes with ABSOLUTELY NO WARRANTY.

pkidb is distributed under the Terms of the GNU General
Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)

Build with go version: %s
`

const helpText = `Usage: %s [--help] [--version] --mode=job|node --down|--drained|--fail|--idle|--up --fini|--time
    --help  Show this help text


    --mode=job|node Trigger mode
                    job  - This programm will be called from job related SLURM triggers
                    node - This programm will be called from node related SLURM triggers

    --retain        MQTT message will be retained by the MQTT broker

    --version       Show version information

  Options for job mode, only one option can be used:
    --down      Node entered DOWN state
    --drained   Node entered DRAINED state
    --fail      Node entered FAIL state
    --idle      Node was idle for a specified amount of time
    --up        Node returned from DOWN state

  Options for node mode, only one option can be used:
    --fini      Job finished
    --time      Job time limit reached

`

const (
	// ModeJob - job mode
	ModeJob int = iota
	// ModeNode - node mode
	ModeNode
)
