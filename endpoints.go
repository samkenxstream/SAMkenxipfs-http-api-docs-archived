// Package docs can be used to gather go-ipfs commands and automatically
// generate documentation or tests.
package docs

import (
	"fmt"
	"sort"

	jsondoc "github.com/Stebalien/go-json-doc"
	cid "github.com/ipfs/go-cid"
	config "github.com/ipfs/go-ipfs"
	cmds "github.com/ipfs/go-ipfs-cmds"
	corecmds "github.com/ipfs/go-ipfs/core/commands"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var JsondocGlossary = jsondoc.NewGlossary().
	WithSchema(new(cid.Cid), jsondoc.Object{"/": "<cid-string>"}).
	WithName(new(multiaddr.Multiaddr), "multiaddr-string").
	WithName(new(peer.ID), "peer-id").
	WithSchema(new(peerstore.PeerInfo),
		jsondoc.Object{"ID": "peer-id", "Addrs": []string{"<multiaddr-string>"}})

var clientOpts = map[string]struct{}{
	cmds.RecLong:     struct{}{},
	cmds.DerefLong:   struct{}{},
	cmds.StdinName:   struct{}{},
	cmds.Hidden:      struct{}{},
	cmds.Ignore:      struct{}{},
	cmds.IgnoreRules: struct{}{},
}

// A map of single endpoints to be skipped or not (subcommands are processed though).
var IgnoreEndpoints = map[string]bool{
	"/api/v0/add":                      false,
	"/api/v0/bitswap/ledger":           true,
	"/api/v0/bitswap/reprovide":        true,
	"/api/v0/bitswap/stat":             true,
	"/api/v0/bitswap/wantlist":         true,
	"/api/v0/block/get":                false,
	"/api/v0/block/put":                false,
	"/api/v0/block/rm":                 true,
	"/api/v0/block/stat":               false,
	"/api/v0/bootstrap":                true,
	"/api/v0/bootstrap/add":            true,
	"/api/v0/bootstrap/add/default":    true,
	"/api/v0/bootstrap/list":           true,
	"/api/v0/bootstrap/rm":             true,
	"/api/v0/bootstrap/rm/all":         true,
	"/api/v0/cat":                      false,
	"/api/v0/cid/base32":               true,
	"/api/v0/cid/bases":                true,
	"/api/v0/cid/codecs":               true,
	"/api/v0/cid/format":               true,
	"/api/v0/cid/hashes":               true,
	"/api/v0/commands":                 true,
	"/api/v0/config":                   true,
	"/api/v0/config/edit":              true,
	"/api/v0/config/profile/apply":     true,
	"/api/v0/config/replace":           true,
	"/api/v0/config/show":              true,
	"/api/v0/dag/export":               true,
	"/api/v0/dag/get":                  false,
	"/api/v0/dag/import":               true,
	"/api/v0/dag/put":                  false,
	"/api/v0/dag/resolve":              false,
	"/api/v0/dht/findpeer":             true,
	"/api/v0/dht/findprovs":            true,
	"/api/v0/dht/get":                  true,
	"/api/v0/dht/provide":              true,
	"/api/v0/dht/put":                  true,
	"/api/v0/dht/query":                true,
	"/api/v0/diag/cmds":                true,
	"/api/v0/diag/cmds/clear":          true,
	"/api/v0/diag/cmds/set-time":       true,
	"/api/v0/diag/sys":                 true,
	"/api/v0/dns":                      true,
	"/api/v0/file/ls":                  true,
	"/api/v0/files/chcid":              true,
	"/api/v0/files/cp":                 true,
	"/api/v0/files/flush":              true,
	"/api/v0/files/ls":                 true,
	"/api/v0/files/mkdir":              true,
	"/api/v0/files/mv":                 true,
	"/api/v0/files/read":               true,
	"/api/v0/files/rm":                 true,
	"/api/v0/files/stat":               true,
	"/api/v0/files/write":              true,
	"/api/v0/filestore/dups":           true,
	"/api/v0/filestore/ls":             true,
	"/api/v0/filestore/verify":         true,
	"/api/v0/get":                      false,
	"/api/v0/id":                       true,
	"/api/v0/key/gen":                  true,
	"/api/v0/key/list":                 true,
	"/api/v0/key/rename":               true,
	"/api/v0/key/rm":                   true,
	"/api/v0/log/level":                true,
	"/api/v0/log/ls":                   true,
	"/api/v0/log/tail":                 true,
	"/api/v0/ls":                       true,
	"/api/v0/mount":                    true,
	"/api/v0/name/publish":             true,
	"/api/v0/name/pubsub/cancel":       true,
	"/api/v0/name/pubsub/state":        true,
	"/api/v0/name/pubsub/subs":         true,
	"/api/v0/name/resolve":             true,
	"/api/v0/object/data":              false,
	"/api/v0/object/diff":              true,
	"/api/v0/object/get":               true,
	"/api/v0/object/links":             true,
	"/api/v0/object/new":               true,
	"/api/v0/object/patch/add-link":    true,
	"/api/v0/object/patch/append-data": true,
	"/api/v0/object/patch/rm-link":     true,
	"/api/v0/object/patch/set-data":    true,
	"/api/v0/object/put":               false,
	"/api/v0/object/stat":              false,
	"/api/v0/p2p/close":                true,
	"/api/v0/p2p/forward":              true,
	"/api/v0/p2p/listen":               true,
	"/api/v0/p2p/ls":                   true,
	"/api/v0/p2p/stream/close":         true,
	"/api/v0/p2p/stream/ls":            true,
	"/api/v0/pin/add":                  false,
	"/api/v0/pin/ls":                   true,
	"/api/v0/pin/rm":                   true,
	"/api/v0/pin/update":               true,
	"/api/v0/pin/verify":               true,
	"/api/v0/ping":                     true,
	"/api/v0/pubsub/ls":                true,
	"/api/v0/pubsub/peers":             true,
	"/api/v0/pubsub/pub":               true,
	"/api/v0/pubsub/sub":               true,
	"/api/v0/refs":                     true,
	"/api/v0/refs/local":               true,
	"/api/v0/repo/fsck":                true,
	"/api/v0/repo/gc":                  true,
	"/api/v0/repo/stat":                true,
	"/api/v0/repo/verify":              true,
	"/api/v0/repo/version":             true,
	"/api/v0/resolve":                  true,
	"/api/v0/shutdown":                 true,
	"/api/v0/stats/bitswap":            true,
	"/api/v0/stats/bw":                 true,
	"/api/v0/stats/dht":                true,
	"/api/v0/stats/repo":               true,
	"/api/v0/swarm/addrs":              true,
	"/api/v0/swarm/addrs/listen":       true,
	"/api/v0/swarm/addrs/local":        true,
	"/api/v0/swarm/connect":            true,
	"/api/v0/swarm/disconnect":         true,
	"/api/v0/swarm/filters":            true,
	"/api/v0/swarm/filters/add":        true,
	"/api/v0/swarm/filters/rm":         true,
	"/api/v0/swarm/peers":              true,
	"/api/v0/tar/add":                  true,
	"/api/v0/tar/cat":                  true,
	"/api/v0/update":                   true,
	"/api/v0/urlstore/add":             true,
	"/api/v0/version":                  false,
	"/api/v0/version/deps":             true,
}

// How much to indent when generating the response schemas
const IndentLevel = 4

// Failsafe when traversing objects containing objects of the same type
const MaxIndent = 20

// Endpoint defines an IPFS RPC API endpoint.
type Endpoint struct {
	Name        string
	Arguments   []*Argument
	Options     []*Argument
	Description string
	Response    string
	Group       string
}

// Argument defines an IPFS RPC API endpoint argument.
type Argument struct {
	Endpoint    string
	Name        string
	Description string
	Type        string
	Required    bool
	Default     string
}

type sorter []*Endpoint

func (a sorter) Len() int           { return len(a) }
func (a sorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

const APIPrefix = "/api/v0"

// AllEndpoints gathers all the endpoints from go-ipfs.
func AllEndpoints() []*Endpoint {
	return Endpoints(APIPrefix, corecmds.Root)
}

func IPFSVersion() string {
	return config.CurrentVersionNumber
}

// Endpoints receives a name and a go-ipfs command and returns the endpoints it
// defines] (sorted). It does this by recursively gathering endpoints defined by
// subcommands. Thus, calling it with the core command Root generates all
// the endpoints.
func Endpoints(name string, cmd *cmds.Command) (endpoints []*Endpoint) {
	var arguments []*Argument
	var options []*Argument

	ignore := cmd.Run == nil || IgnoreEndpoints[name]
	if !ignore { // Extract arguments, options...
		for _, arg := range cmd.Arguments {
			argType := "string"
			if arg.Type == cmds.ArgFile {
				argType = "file"
			}
			arguments = append(arguments, &Argument{
				Endpoint:    name,
				Name:        arg.Name,
				Type:        argType,
				Required:    arg.Required,
				Description: arg.Description,
			})
		}

		for _, opt := range cmd.Options {
			// skip client-side options
			if _, ok := clientOpts[opt.Names()[0]]; ok {
				continue
			}

			def := fmt.Sprint(opt.Default())
			if def == "<nil>" {
				def = ""
			}
			options = append(options, &Argument{
				Name:        opt.Names()[0],
				Type:        opt.Type().String(),
				Description: opt.Description(),
				Default:     def,
			})
		}

		res := buildResponse(cmd.Type)

		endpoints = []*Endpoint{
			&Endpoint{
				Name:        name,
				Description: cmd.Helptext.Tagline,
				Arguments:   arguments,
				Options:     options,
				Response:    res,
			},
		}
	}

	for n, cmd := range cmd.Subcommands {
		endpoints = append(endpoints,
			Endpoints(fmt.Sprintf("%s/%s", name, n), cmd)...)
	}
	sort.Sort(sorter(endpoints))
	return endpoints
}

func buildResponse(res interface{}) string {
	// Commands with a nil type return text. This is a bad thing.
	if res == nil {
		return "This endpoint returns a `text/plain` response body."
	}
	desc, err := JsondocGlossary.Describe(res)
	if err != nil {
		panic(err)
	}
	return desc
}
