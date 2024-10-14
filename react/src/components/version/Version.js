import axios from 'axios';
import {useEffect, useState, useMemo} from "react";
import { DEV_URL } from '../../shared-resources/constants';
import "./Version.css"

const Version = () => {

    const [version, setVersion] = useState('');

    const url = useMemo(() => (process.env.REACT_APP_API_URL === null || process.env.REACT_APP_API_URL === undefined ? DEV_URL : process.env.REACT_APP_API_URL),[]);


    useEffect(() => {
        
        async function GetData(){
            try{
                const response = await axios.get(`${url}/v1/version`)

                setVersion(response?.data?.version)
            } catch {
                setVersion("X.X")
            }
        }

        GetData();

    },[url]);
    return <div class="version-style">
        v{version}
    </div>

}

export default Version;