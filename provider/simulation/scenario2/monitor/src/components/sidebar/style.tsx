import React, { useState, useEffect } from "react";
import theme from "../../styles/theme";
import { styled } from "@material-ui/core/styles";
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
import StarBorder from "@material-ui/icons/StarBorder";
import ListItemIcon from "@material-ui/core/ListItemIcon";

export const ListTitle = styled(({ children }) => (
    <ListItem>
        <Typography>{children}</Typography>
    </ListItem>
))({
    width: "100%"
});

export const ListButton = styled(({ children, ...props }) => (
    <ListItem {...props} button>
        <ListItemText>{children}</ListItemText>
    </ListItem>
))({
    width: "100%"
});

export const CollapseList = styled(
    ({ open, onClickOpen, text, children, ...props }) => (
        <div>
            <ListItem {...props} button onClick={() => onClickOpen()}>
                <ListItemText>{text}</ListItemText>
                {open ? <ExpandLess /> : <ExpandMore />}
            </ListItem>
            <Collapse in={open} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    {children}
                </List>
            </Collapse>
        </div>
    )
)({
    width: "100%"
});

export const SideButton = styled(Button)({
    width: "100%"
});
