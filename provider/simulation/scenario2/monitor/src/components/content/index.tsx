import React, { useState, useEffect } from "react";
import { Typography, Grid, Paper } from "@material-ui/core";
import { Provider, Log, ContentType } from "../../types";
import { styled } from "@material-ui/core/styles";
import AreaMonitor from "./area-monitor";
import LogMonitor from "./log-monitor";

interface Props {
    providers: Provider[];
    contentType: ContentType;
}

const drawerWidth = 240;
const headerHeight = 70;

const ContentContainer = styled("div")({
    paddingTop: headerHeight + 20,
    paddingRight: 20,
    paddingLeft: drawerWidth + 20
});

const showLogs = (provider: Provider) => {
    //console.log("provider", provider);
    return provider
        .getLogs()
        .map((log: Log) => <LogText>{log.Description}</LogText>);
};

const LogGrid = styled(Grid)({
    height: 500,
    margin: 10
});

const LogText = styled(Typography)({
    textAlign: "left"
});

const LogTitle = styled(Typography)({
    margin: 20
});

const LogPaper = styled(Paper)({
    height: 500
});

const LogContent = styled("div")({
    height: 400,
    margin: 10,
    overflow: "auto"
});

const Content: React.FC<Props> = props => {
    const { providers, contentType } = props;

    if (contentType === ContentType.AREA_MONITOR) {
        return <AreaMonitor providers={providers} />;
    } else {
        return <LogMonitor providers={providers} />;
    }
    /*return (
        <ContentContainer>
            <Grid container spacing={4}>
                {providers.map((provider: Provider) => (
                    <LogGrid item xl={4} lg={4} md={6} sm={12} xs={12}>
                        <LogPaper>
                            <LogTitle variant="h6">
                                {provider.getName()}
                            </LogTitle>
                            <LogContent>{showLogs(provider)}</LogContent>
                        </LogPaper>
                    </LogGrid>
                ))}
            </Grid>
        </ContentContainer>
    );*/
};

export default Content;
