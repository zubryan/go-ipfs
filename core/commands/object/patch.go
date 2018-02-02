package objectcmd

import (
	"io"
	"strings"

	cmds "github.com/ipfs/go-ipfs/commands"
	e "github.com/ipfs/go-ipfs/core/commands/e"

	cmdkit "gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
)

var ObjectPatchCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Create a new merkledag object based on an existing one.",
		ShortDescription: `
'ipfs object patch <root> <cmd> <args>' is a plumbing command used to
build custom DAG objects. It mutates objects, creating new objects as a
result. This is the Merkle-DAG version of modifying an object.
`,
	},
	Arguments: []cmdkit.Argument{},
	Subcommands: map[string]*cmds.Command{
		"append-data": patchAppendDataCmd,
		"add-link":    patchAddLinkCmd,
		"rm-link":     patchRmLinkCmd,
		"set-data":    patchSetDataCmd,
	},
}

func objectMarshaler(res cmds.Response) (io.Reader, error) {
	v, err := unwrapOutput(res.Output())
	if err != nil {
		return nil, err
	}

	o, ok := v.(*Object)
	if !ok {
		return nil, e.TypeErr(o, v)
	}

	return strings.NewReader(o.Hash + "\n"), nil
}

var patchAppendDataCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Append data to the data segment of a dag node.",
		ShortDescription: `
Append data to what already exists in the data segment in the given object.

Example:

	$ echo "hello" | ipfs object patch $HASH append-data

NOTE: This does not append data to a file - it modifies the actual raw
data within an object. Objects have a max size of 1MB and objects larger than
the limit will not be respected by the network.
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("root", true, false, "The hash of the node to modify."),
		cmdkit.FileArg("data", true, false, "Data to append.").EnableStdin(),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		api, err := req.InvocContext().GetApi()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		root, err := api.ParsePath(req.Context(), req.StringArguments()[0], api.WithResolve(true))
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		data, err := req.Files().NextFile()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		p, err := api.Object().AppendData(req.Context(), root, data)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(&Object{Hash: p.Cid().String()})
	},
	Type: Object{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: objectMarshaler,
	},
}

var patchSetDataCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Set the data field of an IPFS object.",
		ShortDescription: `
Set the data of an IPFS object from stdin or with the contents of a file.

Example:

    $ echo "my data" | ipfs object patch $MYHASH set-data
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("root", true, false, "The hash of the node to modify."),
		cmdkit.FileArg("data", true, false, "The data to set the object to.").EnableStdin(),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		api, err := req.InvocContext().GetApi()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
		root, err := api.ParsePath(req.Context(), req.StringArguments()[0], api.WithResolve(true))
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		data, err := req.Files().NextFile()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		p, err := api.Object().SetData(req.Context(), root, data)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(&Object{Hash: p.Cid().String()})
	},
	Type: Object{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: objectMarshaler,
	},
}

var patchRmLinkCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Remove a link from an object.",
		ShortDescription: `
Removes a link by the given name from root.
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("root", true, false, "The hash of the node to modify."),
		cmdkit.StringArg("link", true, false, "Name of the link to remove."),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		api, err := req.InvocContext().GetApi()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		root, err := api.ParsePath(req.Context(), req.Arguments()[0], api.WithResolve(true))
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		link := req.Arguments()[1]
		p, err := api.Object().RmLink(req.Context(), root, link)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(&Object{Hash: p.Cid().String()})
	},
	Type: Object{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: objectMarshaler,
	},
}

var patchAddLinkCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Add a link to a given object.",
		ShortDescription: `
Add a Merkle-link to the given object and return the hash of the result.

Example:

    $ EMPTY_DIR=$(ipfs object new unixfs-dir)
    $ BAR=$(echo "bar" | ipfs add -q)
    $ ipfs object patch $EMPTY_DIR add-link foo $BAR

This takes an empty directory, and adds a link named 'foo' under it, pointing
to a file containing 'bar', and returns the hash of the new object.
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("root", true, false, "The hash of the node to modify."),
		cmdkit.StringArg("name", true, false, "Name of link to create."),
		cmdkit.StringArg("ref", true, false, "IPFS object to add link to."),
	},
	Options: []cmdkit.Option{
		cmdkit.BoolOption("create", "p", "Create intermediary nodes."),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		api, err := req.InvocContext().GetApi()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		root, err := api.ParsePath(req.Context(), req.Arguments()[0], api.WithResolve(true))
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		name := req.Arguments()[1]

		child, err := api.ParsePath(req.Context(), req.Arguments()[2], api.WithResolve(true))
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		create, _, err := req.Option("create").Bool()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		p, err := api.Object().AddLink(req.Context(), root, name, child,
			api.Object().WithCreate(create))
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(&Object{Hash: p.Cid().String()})
	},
	Type: Object{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: objectMarshaler,
	},
}
