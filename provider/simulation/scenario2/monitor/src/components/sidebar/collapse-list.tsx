import React, { useState, useEffect } from "react";
import {
    List,
    ListItem,
    Button,
    Typography,
    ListItemText
} from "@material-ui/core";
import ExpandLess from "@material-ui/icons/ExpandLess";
import ExpandMore from "@material-ui/icons/ExpandMore";
import Collapse from "@material-ui/core/Collapse";

interface Props {
    children?: any;
    title: string;
}

const CollapseList: React.FC<Props> = (props: Props) => {
    const { children, title } = props;
    const [open, setOpen] = useState<boolean>(false);

    const handleOpen = () => {
        setOpen(!open);
    };

    return (
        <div>
            <ListItem disableGutters button onClick={() => handleOpen()}>
                <ListItemText>{title}</ListItemText>
                {open ? <ExpandLess /> : <ExpandMore />}
            </ListItem>
            <Collapse in={open} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    {children}
                </List>
            </Collapse>
        </div>
    );
};

export default CollapseList;
