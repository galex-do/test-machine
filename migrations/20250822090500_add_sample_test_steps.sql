-- +goose Up
-- Insert sample test steps data for all test cases
INSERT INTO test_steps (test_case_id, step_number, description, expected_result) VALUES
    -- Test Case 1: Valid User Login
    (1, 1, 'Navigate to the login page', 'Login page loads successfully with username and password fields'),
    (1, 2, 'Enter valid username and password', 'Credentials are accepted and fields show no validation errors'),
    (1, 3, 'Click the login button', 'User is redirected to dashboard and sees welcome message'),
    
    -- Test Case 2: Invalid Password Login
    (2, 1, 'Navigate to the login page', 'Login page loads successfully'),
    (2, 2, 'Enter valid username but invalid password', 'Password field shows validation error'),
    (2, 3, 'Click the login button', 'Error message displays: "Invalid password"'),
    
    -- Test Case 3: Password Reset Flow
    (3, 1, 'Navigate to login page and click "Forgot Password"', 'Password reset form is displayed'),
    (3, 2, 'Enter valid email address', 'Email field accepts input and shows no validation errors'),
    (3, 3, 'Click submit button', 'Success message shows: "Password reset link sent to your email"'),
    
    -- Test Case 4: Add Item to Cart
    (4, 1, 'Browse to product catalog page', 'Product catalog loads with available items'),
    (4, 2, 'Select a product and click "Add to Cart"', 'Product is added to cart and cart counter increases'),
    
    -- Test Case 5: Checkout Process
    (5, 1, 'Add items to cart and navigate to checkout', 'Checkout page displays with cart items and total'),
    (5, 2, 'Fill in shipping and billing information', 'Forms accept input and validate required fields'),
    (5, 3, 'Review order and click "Place Order"', 'Order confirmation page displays with order number'),
    
    -- Test Case 6: Payment Processing
    (6, 1, 'Proceed to payment section during checkout', 'Payment form loads with supported payment methods'),
    (6, 2, 'Enter valid credit card information', 'Payment form validates and accepts card details'),
    (6, 3, 'Submit payment', 'Payment processes successfully and confirmation is displayed'),
    
    -- Test Case 7: Navigation Menu
    (7, 1, 'Open the mobile application', 'Application launches successfully and shows the main screen'),
    (7, 2, 'Tap on the navigation menu button', 'Navigation menu slides out from the left side'),
    (7, 3, 'Verify all menu items are visible', 'All navigation options are displayed: Home, Profile, Settings, Logout'),
    
    -- Test Case 8: Screen Transitions
    (8, 1, 'Navigate from home screen to profile', 'Screen transitions smoothly with appropriate animation'),
    (8, 2, 'Return to home screen using back button', 'Transition back is smooth and maintains app state')
ON CONFLICT (test_case_id, step_number) DO NOTHING;

-- +goose Down
DELETE FROM test_steps WHERE test_case_id IN (1, 2, 3, 4, 5, 6, 7, 8);