import axios from 'axios';

export const fetchAPTData = async () => {
  const response = await axios.get(`/api/APT`);
  return response.data;
};

export const fetchRemndrAPTData = async () => {
  const response = await axios.get(`/api/RemndrAPT`);
  return response.data;
};