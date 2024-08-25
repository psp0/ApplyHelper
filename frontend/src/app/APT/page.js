'use client'
import React, { useEffect, useState } from 'react';

import SearchForm from './SearchForm';
import Table from './Table';
import Modal from './Modal';
import Menu from '../Menu';
import './page.style.css'; 
import '../page.style.css'; 
import '../spinner.style.css'
import { fetchAPTData } from '../../../lib/api';
import { fetchdevAPTData } from '../../../lib/dev-api';

const APT = () => {
  const [aptData, setAptData] = useState({ data: [] });
  const [modalData, setModalData] = useState(null);
  const [houseName, setHouseName] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(true); 
  const [isFetching, setIsFetching] = useState(false); 

  useEffect(() => {
    const loadData = async () => {
      try {
        let data;
        if (process.env.NODE_ENV === 'development') {
          data = await fetchdevAPTData(currentPage);
        } else {
          data = await fetchAPTData(currentPage); 
        }
        
        if (data && data.data && data.data.length > 0) {
          setAptData((prev) => ({
            data: [...prev.data, ...data.data], 
          }));
        } else {
          // Stop fetching more data if no data is returned
          setIsFetching(false);
          setLoading(false); // Add this line to stop the loading spinner
          return;
        }
        
        setLoading(false);
        setIsFetching(false); 
      } catch (error) {
        alert('Failed to fetch data');
        setLoading(false);
        setIsFetching(false);
      }
    };

    loadData();
  }, [currentPage]);

  useEffect(() => {
    const handleScroll = () => {
      if (window.innerHeight + document.documentElement.scrollTop + 100 >= document.documentElement.scrollHeight && !isFetching) {
        setIsFetching(true);
        setCurrentPage((prevPage) => prevPage + 1);
      }
    };

    window.addEventListener('scroll', handleScroll);

    return () => window.removeEventListener('scroll', handleScroll);
  }, [isFetching]);

  const handleButtonClick = (detailData, houseDtlSecdNm, houseName) => {
    setModalData(detailData);
    setHouseName(houseName);
  };

  const closeModal = () => {
    setModalData(null);
  };

  const loadMore = () => {
    setIsFetching(true);
    setCurrentPage((prevPage) => prevPage + 1);
  };

  return (
    <div>
      <header>
        <h1>청약지기</h1>
      </header>
      <Menu/>
      <SearchForm setAptData={setAptData} setLoading={setLoading} />
      {loading ? ( 
        <div className="spinner"></div>
      ) : (
        <>
          <Table aptData={aptData.data} handleButtonClick={handleButtonClick} />
          {isFetching ? (
            <div className="spinner"></div> 
          ) : (
            <button onClick={loadMore} className="load-more-btn">
              더보기
            </button>
          )}
        </>
      )}
      {modalData && <Modal modalData={modalData}  houseName={houseName} closeModal={closeModal} />}
    </div>
  );
};

export default APT;
