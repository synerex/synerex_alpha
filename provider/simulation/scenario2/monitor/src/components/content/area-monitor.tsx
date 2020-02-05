import React, { useState, useEffect } from "react";
import { Typography, Grid, Paper } from "@material-ui/core";
import { Provider, Log } from "../../types";
import { styled } from "@material-ui/core/styles";
import GoogleMapReact from "google-map-react";

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
const AnyReactComponent = (props: any) => <div>{props.text}</div>;

const AreaMonitor: React.FC<Props> = props => {
    const { providers } = props;
    const defaultProps = {
        center: {
            lat: 59.95,
            lng: 30.33
        },
        zoom: 11
    };
    return (
        <ContentContainer>
            <Typography variant="h2">{"Area Monitor"}</Typography>
            <div style={{ height: "100vh", width: "100%" }}>
                <GoogleMapReact
                    bootstrapURLKeys={{
                        key: "AIzaSyBj9mm1Y-7mnZx2Vh1DLrLhWwZ9taRqAI0"
                    }}
                    defaultCenter={defaultProps.center}
                    defaultZoom={defaultProps.zoom}
                >
                    <AnyReactComponent
                        lat={59.955413}
                        lng={30.337844}
                        text="My Marker"
                    />
                </GoogleMapReact>
            </div>
        </ContentContainer>
    );
};

export default AreaMonitor;
