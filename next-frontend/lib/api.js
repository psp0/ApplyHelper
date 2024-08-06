// lib/api.js
import axios from 'axios';

export const fetchMainData = async () => {
  const response = await axios.get(`http://localhost:8080/data`);
  return response.data;
};
