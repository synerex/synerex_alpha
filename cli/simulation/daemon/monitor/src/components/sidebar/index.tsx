import React, { useState, useEffect } from "react";
import {
    Divider,
    Drawer,
    createStyles,
    Theme,
    Typography,
    TextField
} from "@material-ui/core";
import theme from "../../styles/theme";
import { List, ListItem, Button } from "@material-ui/core";
import { Provider, Command, Option } from "../../types";
import { SideButton, ListButton, ListTitle } from "./style";
import ListItemText from "@material-ui/core/ListItemText";
import Collapse from "@material-ui/core/Collapse";
import InboxIcon from "@material-ui/icons/MoveToInbox";
import CollapseList from "./collapse-list";
import OptionLists from "./option-lists";

interface Props {
    providers: Provider[];
    commands: Command[];
    runCommand: (command: Command) => void;
}

const drawerWidth = 240;
const headerHeight = 70;

const Sidebar: React.FC<Props> = props => {
    const { providers, commands, runCommand } = props;

    return (
        <Drawer
            anchor="left"
            variant={"permanent"}
            style={{
                flexShrink: 0
            }}
        >
            <div style={{ width: drawerWidth, marginTop: headerHeight }}>
                {/*<List dense={false}>
                    <ListTitle>{"Providers"}</ListTitle>

                    {providers.map((provider: Provider, i) => (
                        <CollapseList
                            //key={provider.Name}
                            title={provider.getName()}
                        >
                            <OptionLists
                                target={provider}
                                handleRun={(provider: Provider) =>
                                    runProvider(provider)
                                }
                            />
                        </CollapseList>
                    ))}
							</List>*/}

                <List>
                    <ListTitle>{"Commands"}</ListTitle>
                    {commands.map((command, i) => (
                        <CollapseList
                            //key={command.Name}
                            title={command.Name}
                        >
                            <OptionLists
                                command={command}
                                runCommand={(command: Command) =>
                                    runCommand(command)
                                }
                            />
                        </CollapseList>
                    ))}
                </List>
            </div>
        </Drawer>
    );
};

export default Sidebar;
