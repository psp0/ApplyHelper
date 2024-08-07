import axios from 'axios';

export const fetchAPTData = async () => {
  const response = await axios.get(`http://localhost:8080/getAPT`);
  return response.data;
};


export const fetchRemndrAPTData = async () => {
  const response = await axios.get(`http://localhost:8080/getRemndrAPT`);
  return response.data;
};