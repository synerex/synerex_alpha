import React, { useState, useEffect } from "react";
import { Divider, Drawer } from "@material-ui/core";
import theme from "../../styles/theme";
import { List, ListItem, Button, Modal } from "@material-ui/core";
import { Provider, Command, Option, DialogStatus } from "../../types";
import { SideButton, ListButton, ListTitle } from "./style";
import ListItemText from "@material-ui/core/ListItemText";
import CommandDialog from "./dialog";

interface Props {
    providers: Provider[];
    commands: Command[];
    runCommand: (command: Command) => boolean;
}

const drawerWidth = 240;
const headerHeight = 70;

const Sidebar: React.FC<Props> = props => {
    const { providers, commands, runCommand } = props;

    const [dialogStatus, setDialogStatus] = useState<DialogStatus>({
        Open: false
    });

    const handleOpen = (command: Command) => {
        setDialogStatus({ Open: true, Command: command });
    };

    const handleClose = () => {
        setDialogStatus({ Open: false });
    };

    return (
        <Drawer
            anchor="left"
            variant={"permanent"}
            style={{
                flexShrink: 0
            }}
        >
            <div style={{ width: drawerWidth, marginTop: headerHeight }}>
                <List>
                    <ListTitle>{"Commands"}</ListTitle>
                    <Divider />
                    {commands.map((command, i) => (
                        <div>
                            {" "}
                            <ListItem
                                button
                                onClick={() => handleOpen(command)}
                            >
                                <ListItemText primary={command.Name} />
                            </ListItem>
                        </div>
                    ))}
                </List>
                <CommandDialog
                    dialogStatus={dialogStatus}
                    handleClose={handleClose}
                    handleRun={runCommand}
                />
            </div>
        </Drawer>
    );
};

export default Sidebar;
