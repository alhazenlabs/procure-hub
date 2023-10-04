// AppRouter.js
import React from 'react';
import { Routes, Route } from 'react-router-dom';
import SignIn from './SignIn'; // Import your Signin
import CataloguePage from './CataloguePage'; // Import your catalogue component

const AppRouter = () => {
  return (
    <Routes>
        <Route path="/" element={<SignIn />} />
        <Route path="/catalogue" component={<CataloguePage />} />
        {/* Add more routes as needed */}
    </Routes>
  );
};

export default AppRouter;