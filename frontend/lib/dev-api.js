import axios from 'axios';

export const fetchdevAPTData = async (page) => {
  const response = await axios.get(`${process.env.BACKEND_URL}/APT?page=${page}`);  
  return response.data;
};

export const fetchdevRemndrAPTData = async (page) => {
  const response = await axios.get(`${process.env.BACKEND_URL}/RemndrAPT?page=${page}`);
  return response.data;
};