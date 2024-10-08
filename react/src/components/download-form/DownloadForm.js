import {useState} from 'react'
import axios from 'axios';

const DownloadForm = () => {

    const [formData, setFormData] = useState({url:'',actionType:'p'})

    async function handleForm(event){
        event.preventDefault();

        try{
            const response = await axios.post('http://localhost:8080/api/v1/music',{
                url: formData.url,
                actionType: formData.actionType
            },{
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
            console.log(err)
            alert(await err?.response?.data?.text())
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


    return <div>
        <form onSubmit={handleForm}>
            <label>Enter URL
                <input type="text" name="url"placeholder="Past URL here" value={formData.url} onChange={handleChange} required/>
            </label>
            <select name="actionType" value={formData.actionType} onChange={handleChange}>
                <option value="p">Play List</option>
                <option value="t">Single Track</option>
                <option value="a">Album</option>
            </select>
            <button type="submit">Download</button>
        </form>
    </div>
}

export default DownloadForm;