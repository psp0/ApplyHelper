'use client'
import React, { useEffect, useState } from 'react';
import { fetchMainData } from '../../lib/api';
import './page.style.css'; 
const Home = () => {
  const [aptData, setAptData] = useState({ data: [] });
  const [modalData, setModalData] = useState(null);
  const [isKukMin, setIsKukMin] = useState(false);
  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await fetchMainData();
        console.log(data);
        setAptData(data);
      } catch (error) {
        alert('Failed to fetch data');
      }
    };

    loadData();
  }, []);

  const handleButtonClick = (detailData,houseDtlSecdNm) => {
    if (houseDtlSecdNm==="국민") setIsKukMin(true)
    setModalData(detailData);    
  };

  const closeModal = () => {
    setModalData(null);
    setIsKukMin(false)
  };

  return (
    <div>
      <table className="tbl_st">
        <thead>
          <tr>
          <th>순서</th>
        
            <th>지역</th>
            <th>주택 구분</th>            
            <th>주택명</th> 
            <th>모집공고일</th>  
            <th>일반공급 물량 및 경쟁률</th>
          </tr>
        </thead>
        <tbody>
          {aptData.data && aptData.data.length > 0 ? (
            aptData.data.map((e,i) => (
              <tr key={e.HOUSE_MANAGE_NO}>
                <td>{i}</td>                
                <td>{e.HOUSE_DTL_SECD_NM}</td>
                <td>{e.SUBSCRPT_AREA_CODE_NM}</td>                     
                <td>{e.HOUSE_NM}</td>         
                <td>{e.RCRIT_PBLANC_DE}</td>
                <td>
                  <button type="button" onClick={() => handleButtonClick(e.DETAIL_DATA,e.HOUSE_DTL_SECD_NM)}>물량</button>
                </td>
              </tr>
            ))
          ) : (
            <tr><td colSpan="6">데이터가 없습니다</td></tr>
          )}
        </tbody>
      </table>

      {modalData && (
        <div className="modal">
          <div className="modal-content">
            <span className="close" onClick={closeModal}>&times;</span>
     
                {modalData.map((item, index) => {
                  const localData = item.LOCAL_POINT || item.LOCAL_RAND_ZERO || item.LOCAL_RAND_ZERO_ONE;
                  const etcGGData = item.ETC_GG_POINT || item.ETC_GG_RAND_ZERO || item.ETC_GG_RAND_ZERO_ONE;
                  const etcData = item.ETC_POINT || item.ETC_RAND_ZERO || item.ETC_RAND_ZERO_ONE;
                  if (isKukMin){
                    return (
                      <table>
                      <tbody>
                      <tr>
                                <th >타입</th>
                                <th >가격</th>                                 
                                    <th >지역</th>
                                    <th >구분</th>
                                    <th >총물량</th>
                                  </tr>
                          <td rowSpan="0">{item.HOUSE_TY}</td>         
                          <td rowSpan="0">{item.LTTOT_TOP_AMOUNT}억</td>                           
                            <td rowSpan="3">해당</td>
                            <td>{Number(item.HOUSE_TY.replace(/[^0-9.]/g, '')) > 40? "高저축총액자":"高납입횟수자"}</td>
                            <td>{item.SUPLY_HSHLDCO}</td>                         
                                  </tbody>
                                  </table>
                    )                    
                  }
                  return (
                    <table>
                         
              <tbody>          
              <tr>
                        <th >타입</th>
                        <th >가격</th>
                            <th >총물량</th>
                            <th >지역</th>
                            <th >구분</th>
                            <th >물량</th>
                          </tr>

                    <React.Fragment key={index}>
                      {localData>0  && (
                        <> 
                     
                          <tr>        
                          <td rowSpan="0">{item.HOUSE_TY}</td>         
                          <td rowSpan="0">{item.LTTOT_TOP_AMOUNT}억</td>              
                          <td rowSpan="0">{item.SUPLY_HSHLDCO}</td>           
                            <td rowSpan="3">해당</td>
                            <td>高가점자</td>
                            <td>{item.LOCAL_POINT}</td>
                          </tr>
                          <tr>
                            <td>랜덤 - 무주택</td>
                            <td>{item.LOCAL_RAND_ZERO}</td>
                          </tr>
                          <tr>
                            <td>랜덤 - 무주택,일주택</td>
                            <td>{item.LOCAL_RAND_ZERO_ONE}</td>
                          </tr>
                        </>
                      )}
                      { etcGGData >0  && (
                        <>
                          <tr>
                            <td rowSpan="3">해당, 기타경기</td>
                            <td>高가점자</td>
                            <td>{item.ETC_GG_POINT}</td>
                          </tr>
                          <tr>
                            <td>랜덤 - 무주택</td>
                            <td>{item.ETC_GG_RAND_ZERO}</td>
                          </tr>
                          <tr>
                            <td>랜덤 - 무주택,일주택</td>
                            <td>{item.ETC_GG_RAND_ZERO_ONE}</td>
                          </tr>
                        </>
                      )}
                      { etcData >0 && (
                        <>
                          <tr>
                            {
                              etcGGData >0 &&   <td rowSpan="3">해당, 기타경기, 기타</td>
                            }
                            {
                              etcGGData <= 0 &&  <td rowSpan="3">해당, 기타</td>
                            }                            
                            <td>高가점자</td>
                            <td>{item.ETC_POINT}</td>
                          </tr>
                          <tr>
                            <td>랜덤 - 무주택</td>
                            <td>{item.ETC_RAND_ZERO}</td>
                          </tr>
                          <tr>
                            <td>랜덤 - 무주택,일주택</td>
                            <td>{item.ETC_RAND_ZERO_ONE}</td>
                          </tr>
                        </>
                      )}
                    </React.Fragment>
                    </tbody>
                    </table>
                  );
                })}

                  <div>
                    <p className="highlight">
                      ※ 특별공급 미달 잔여물량은 포함안함
                      <br />
                      ※ 할당된 물량은 모두 가져간다고 가정
                      <br />
                      ※ 투기 과열지역에 대한 계산은 미구현
                    </p>
                  </div>
            {/* 만약 경쟁률 계산시 이를 고려해야함 */}
            {/* 물량확인 버튼, 날짜 넘어가면 경쟁률 확인 버튼 */}
          </div>
        </div>
      )}
    </div>
  );
};

export default Home;
