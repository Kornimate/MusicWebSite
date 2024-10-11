import {useState, useMemo} from 'react'
import { DEV_URL } from '../../shared-resources/constants';
import axios from 'axios';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Box from '@mui/material/Box';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import logo from "../../logo.svg";
import "./DownloadForm.css"
import { CircularProgress } from '@mui/material';

const DownloadForm = () => {

    const [formData, setFormData] = useState({url:''});
    const [isLoading, setIsLoading] = useState(false);
    const [open,setOpen] = useState(false);
    const [alertText, setAlertText] = useState('');
    const [alertColor, setAlertColor] = useState('success');

    const url = useMemo(() => (process.env.REACT_APP_API_URL === null || process.env.REACT_APP_API_URL === undefined ? DEV_URL : process.env.REACT_APP_API_URL),[]);

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {
          return;
        }
    
        setOpen(false);
    };

    async function handleForm(event){
        event.preventDefault();

        setIsLoading(true);

        try{
            const response = await axios.post(`${url}/v1/music`,formData,{
                responseType: 'blob'
            });

            if(response?.status === 500){
                throw new Error("Request failed");
            }

            let fileName = response.headers["filename"];

            if (fileName === undefined || fileName === null){
                fileName = "MusicApp.zip";
            }

            downloadMusicFile(new Blob([response.data]), fileName);

            setIsLoading(false);
            setAlertText("file successfully downloaded");
            setAlertColor('success');
            setOpen(true);
        } catch (err) {
            let error = await err?.response?.data?.text();

            if(error === undefined)
                error = "Error while downloading file";

            setIsLoading(false);
            setAlertText(error);
            setAlertColor('error')
            setOpen(true);

            console.log(error)
        }
    }

    function handleChange(event){
        const {value, name} = event.target;

        setFormData((prevData) => ({...prevData, [name]:value}))
    }

    function downloadMusicFile(fileBlob, fileName){

        const windowUrl = window.URL.createObjectURL(fileBlob);

        const link = document.createElement("a");

        link.href = windowUrl;

        link.setAttribute("download",fileName);

        document.body.appendChild(link);

        link.click();

        document.body.removeChild(link);

        window.URL.revokeObjectURL(windowUrl);
    }


    return <>
        <Box sx={{display:'block', alignItems:'center', width:'30%'}}>
            <form onSubmit={handleForm} id="form">
                <img src={logo} alt="logo"/>
                { isLoading 
                    ? <Box sx={{display:'block', alignItems:'center', width:'100%'}}>
                        <CircularProgress size="7.7rem" />
                    </Box>
                    : <Box sx={{display:'block', alignItems:'center', width:'100%'}}>
                        <TextField name="url" id="textField" label="Url of video" variant="outlined" value={formData.url} onChange={handleChange} required sx={{width: "100%"}} />
                        <Button variant="outlined" type="submit" size="large" id="submitBtn">Download</Button>
                    </Box>
                }
            </form>
        </Box>
        <Snackbar 
            open={open}
            autoHideDuration={4000}
            onClose={handleClose}
            anchorOrigin={{ vertical:'bottom', horizontal:'right'}}
        >
            <Alert
                onClose={handleClose}
                severity={alertColor}
                variant="filled"
                sx={{ width: '100%' }}
            >
            {alertText}
            </Alert>
    </Snackbar>
  </>
}

export default DownloadForm;