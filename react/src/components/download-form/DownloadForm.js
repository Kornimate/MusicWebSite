import {useState, useMemo} from 'react'
import { DEV_URL } from '../../shared-resources/constants';
import axios from 'axios';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import logo from "../../logo.svg";
import "./DownloadForm.css"

const DownloadForm = () => {

    const [formData, setFormData] = useState({url:''})

    const url = useMemo(() => (process.env.API_URL === null || process.env.API_URL === undefined ? DEV_URL : process.env.API_URL),[]);

    async function handleForm(event){
        event.preventDefault();

        try{
            const response = await axios.post(`${url}/v1/music`,formData,{
                responseType: 'blob'
            });

            if(response?.status === 500){
                throw new Error("Request failed")
            }

            let fileName = response.headers["filename"]

            if (fileName === undefined || fileName === null){
                fileName = "MusicApp.zip"
            }

            console.log(response.headers)
            console.log(fileName)

            downloadMusicFile(new Blob([response.data]), fileName);
        } catch (err) {
            let error = await err?.response?.data?.text()

            if(error === undefined)
                error = "Error while getting data"

            console.log(error)
            alert(error)
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


    return <div id="container">
        <form onSubmit={handleForm} id="form">
            <img src={logo} alt="logo"/>
            <TextField name="url" id="textField" label="Url of video" variant="outlined" value={formData.url} onChange={handleChange} required sx={{width: "100%"}} /><br />
            <Button variant="outlined" type="submit" size="large" id="submitBtn">Download</Button>
        </form>
    </div>
}

export default DownloadForm;