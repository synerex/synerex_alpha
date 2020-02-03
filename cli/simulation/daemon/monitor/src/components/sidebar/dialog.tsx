import React, { useState, useEffect } from "react";
import {
    Button,
    Typography,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    TextField,
    Snackbar
} from "@material-ui/core";
import MuiAlert, { AlertProps } from "@material-ui/lab/Alert";
import { DialogStatus, Command, CommandType } from "../../types";

function Alert(props: AlertProps) {
    return <MuiAlert elevation={6} variant="filled" {...props} />;
}

enum AlertType {
    ERROR,
    SUCCESS,
    NONE
}

interface Props {
    dialogStatus: DialogStatus;
    handleClose: () => void;
    handleRun: (command: Command) => boolean;
}

const CommandDialog: React.FC<Props> = (props: Props) => {
    const { dialogStatus, handleClose, handleRun } = props;

    const [alertType, setAlertType] = useState<AlertType>(AlertType.NONE);

    const handleCloseAlert = () => {
        setAlertType(AlertType.NONE);
    };

    const showSnackbar = () => {
        if (alertType === AlertType.SUCCESS) {
            return (
                <Snackbar
                    open={true}
                    autoHideDuration={6000}
                    onClose={handleCloseAlert}
                >
                    <Alert onClose={handleCloseAlert} severity="success">
                        This is a success message!
                    </Alert>
                </Snackbar>
            );
        } else if (alertType === AlertType.ERROR) {
            return (
                <Snackbar
                    open={true}
                    autoHideDuration={6000}
                    onClose={handleCloseAlert}
                >
                    <Alert onClose={handleCloseAlert} severity="error">
                        This is an error message!
                    </Alert>
                </Snackbar>
            );
        }
    };

    const showContent = () => {
        switch (dialogStatus.Command?.Type) {
            case CommandType.SET_AGENTS:
                return (
                    <DialogContent>
                        <DialogContentText>
                            Input SetAgents Option
                        </DialogContentText>
                        <TextField
                            autoFocus
                            margin="dense"
                            id="type"
                            label="Type"
                            type="text"
                            fullWidth
                        />
                        <TextField
                            autoFocus
                            margin="dense"
                            id="num"
                            label="Num"
                            type="text"
                            fullWidth
                        />
                    </DialogContent>
                );
            case CommandType.SET_AREA:
                return (
                    <DialogContent>
                        <DialogContentText>
                            Input SetArea Option
                        </DialogContentText>
                        <TextField
                            autoFocus
                            margin="dense"
                            id="area"
                            label="AreaCoord"
                            type="text"
                            fullWidth
                        />
                    </DialogContent>
                );
            case CommandType.SET_CLOCK:
                return (
                    <DialogContent>
                        <DialogContentText>
                            Input SetClock Option
                        </DialogContentText>
                        <TextField
                            autoFocus
                            margin="dense"
                            id="time"
                            label="Time"
                            type="text"
                            fullWidth
                        />
                    </DialogContent>
                );
            case CommandType.START_CLOCK:
                return (
                    <DialogContent>
                        <DialogContentText>
                            Input StartClock Option
                        </DialogContentText>
                    </DialogContent>
                );
            case CommandType.STOP_CLOCK:
                return (
                    <DialogContent>
                        <DialogContentText>
                            Input StopClock Option
                        </DialogContentText>
                    </DialogContent>
                );
        }
    };

    return (
        <div>
            <Dialog
                open={dialogStatus.Open}
                onClose={handleClose}
                aria-labelledby="form-dialog-title"
            >
                <DialogTitle id="form-dialog-title">
                    {dialogStatus.Command?.Name}
                </DialogTitle>
                {showContent()}

                <DialogActions>
                    <Button onClick={handleClose} color="primary">
                        Cancel
                    </Button>
                    <Button
                        onClick={() => {
                            if (dialogStatus.Command) {
                                const err = handleRun(dialogStatus.Command);
                                if (err) {
                                    setAlertType(AlertType.ERROR);
                                } else {
                                    setAlertType(AlertType.SUCCESS);
                                    handleClose();
                                }
                            }
                        }}
                        color="primary"
                    >
                        Run
                    </Button>
                </DialogActions>
            </Dialog>
            {showSnackbar()}
        </div>
    );
};

export default CommandDialog;
