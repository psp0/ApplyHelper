'use client'
import React, { useEffect, useState } from 'react';
import Menu from './Menu';
import './page.style.css'; 
import './spinner.style.css'

const Home = () => {
  useEffect(() => {
      
  }, []); 


  return (
    <div>
      <header>
        <h1>청약지기</h1>
      </header>
      <Menu/>
      여기는 홈페이지 입니다.

    </div>
  );
};

export default Home;
