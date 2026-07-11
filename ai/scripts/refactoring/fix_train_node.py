import re

path = "internal/adapter/cli/cli.go"
with open(path, "r") as f: content = f.read()

# Add to runCharacter
content = content.replace('case "give-item":\n\t\treturn a.characterGiveItem(ctx, args[1:])\n\tdefault:',
'''case "give-item":
		return a.characterGiveItem(ctx, args[1:])
	case "train-node":
		return a.characterTrainNode(ctx, args[1:])
	default:''')

content = content.replace('<create|list|give-item>', '<create|list|give-item|train-node>')

# Add the function
train_node_func = '''
func (a *App) characterTrainNode(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character train-node", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.TrainNodeInput
	fs.StringVar(&in.Character, "name", "", "character name")
	fs.StringVar(&in.System, "system", "", "power system name")
	fs.StringVar(&in.NodeID, "node", "", "node ID to train/unlock")
	
	if err := fs.Parse(args); err != nil {
		return 2
	}
	
	if in.Character == "" || in.System == "" || in.NodeID == "" {
		_, _ = fmt.Fprintln(a.err, "name, system, and node are required")
		return 2
	}

	c, err := a.characters.TrainNode(ctx, in)
	if err != nil {
		return a.fail(err)
	}
	
	logger.System("Training complete! %s advanced in %s.", c.Name, in.System)
	_, _ = fmt.Fprintf(a.out, "%s trained node %s in %s. Power is now %s\\n", c.Name, in.NodeID, in.System, c.PowerValue)
	return 0
}
'''

content += train_node_func
with open(path, "w") as f: f.write(content)
