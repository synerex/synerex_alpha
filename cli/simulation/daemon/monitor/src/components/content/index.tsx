import React, { useState, useEffect } from "react";
import { Typography } from "@material-ui/core";
import { Provider } from "../../types";
import { styled } from "@material-ui/core/styles";

interface Props {
    providers: Provider[];
}

const drawerWidth = 240;
const headerHeight = 70;

const ContentContainer = styled("div")({
    paddingTop: headerHeight + 20,
    paddingRight: 20,
    paddingLeft: drawerWidth + 20
});

const Content: React.FC<Props> = props => {
    const { providers } = props;

    return (
        <ContentContainer>
            {providers.map((provider: Provider) => (
                <Typography>{provider.getName()}</Typography>
            ))}
        </ContentContainer>
    );
};

export default Content;
