import React, { useState, useEffect } from "react";
import {
    List,
    ListItem,
    Button,
    ListItemText,
    TextField
} from "@material-ui/core";
import { Option, Provider, Command } from "../../types";

interface Props {
    command: Command;
    runCommand: (command: Command) => void;
}

const OptionLists: React.FC<Props> = props => {
    const { command, runCommand } = props;

    const [option, setOption] = useState<Command["Option"]>(command.Option);

    const handleChange = (value: string, i: number) => {
        const newOptions = [...options];
        newOptions[i] = {
            Key: options[i].Key,
            Value: value
        };
        setOptions(newOptions);
    };

    const handleClick = () => {
        console.log("options: ", options);
        const newCommand: Command = {
            Name: command.Name,
            Options: options
        };
        runCommand(newCommand);
    };
    return (
        <div>
            {Object.keys(option).forEach((k: string, i: number) => (
                <ListItem key={i}>
                    <ListItemText primary={k} />
                    <TextField
                        defaultValue={string(option[k])}
                        onChange={e => handleChange(e.target.value, i)}
                    />
                </ListItem>
            ))}
            <Button onClick={() => handleClick()} variant={"contained"}>
                Run
            </Button>
        </div>
    );
};

export default OptionLists;
