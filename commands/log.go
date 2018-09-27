package commands

import (
	cmds "gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	cmdkit "gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"

	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
)

var logCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Interact with the daemon event log output.",
		ShortDescription: `
'go-filecoin log' contains utility commands to affect the event logging
output of a running daemon.
`,
	},

	Subcommands: map[string]*cmds.Command{
		"tail":   logTailCmd,
		"stream": logStreamCmd,
	},
}

var logTailCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Read the event log.",
		ShortDescription: `
Outputs event log messages (not other log messages) as they are generated.
`,
	},

	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) {
		r := GetAPI(env).Log().Tail(req.Context)
		re.Emit(r) // nolint: errcheck
	},
}

var logStreamCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Stream the event logs to a multiaddress.",
		ShortDescription: `
Outputs event log messages (not other log messages) as they are generated to
 the specified multiaddress
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("addr", true, false, "multiaddress logs will stream to"),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) {
		maddr, err := ma.NewMultiaddr(req.Arguments[0])
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal)
			return
		}
		if err := GetAPI(env).Log().Stream(req.Context, maddr); err != nil {
			re.SetError(err, cmdkit.ErrNormal)
			return
		}
	},
}
