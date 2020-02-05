import React, { useState } from "react";
import {
    Badge,
    AppBar,
    Toolbar,
    Hidden,
    IconButton,
    Button,
    createStyles,
    Theme
} from "@material-ui/core";
import { styled } from "@material-ui/core/styles";
import theme from "../../styles/theme";
import { ContentType } from "../../types";

const LogoImage = styled("img")({
    background: theme.palette.primary.main,
    color: "white",
    height: 48
});

const Title = styled(Button)({
    background: theme.palette.primary.main,
    fontSize: 20,
    textAlign: "start",
    color: "white"
});

const Div = styled("div")({
    background: theme.palette.primary.main,
    //flexGrow: 1
    width: 300
});

const SignUpButton = styled(Button)({
    background: theme.palette.primary.main,
    marginLeft: 10,
    marginRight: 10,
    color: "white",
    height: 48
});

const SignInButton = styled(Button)({
    background: theme.palette.primary.main,
    marginLeft: 10,
    marginRight: 10,
    color: "white",
    height: 48
});

interface Props {
    changeContent: (type: ContentType) => void;
}

const Header: React.FC<Props> = props => {
    const { changeContent } = props;
    return (
        <AppBar position="fixed" style={{ zIndex: theme.zIndex.drawer + 1 }}>
            <Toolbar>
                <Button onClick={() => changeContent(ContentType.LOG_MONITOR)}>
                    {"Log Monitor"}
                </Button>
                <Button onClick={() => changeContent(ContentType.AREA_MONITOR)}>
                    {"Area Monitor"}
                </Button>
            </Toolbar>
        </AppBar>
    );
};

export default Header;
