import axios from 'axios';

export const fetchAPTData = async (page) => {
  const response = await axios.get(`/api/APT?page=${page}`);
  return response.data;
};

export const fetchRemndrAPTData = async (page) => {
  const response = await axios.get(`/api/RemndrAPT?page=${page}`);
  return response.data;
};